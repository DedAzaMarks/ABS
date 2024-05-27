package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/DedAzaMarks/ABS/internal/domain"
	myerrors "github.com/DedAzaMarks/ABS/internal/domain/errors"
	pb "github.com/DedAzaMarks/ABS/internal/proto"
	"github.com/DedAzaMarks/ABS/internal/server/scraper"
	"github.com/DedAzaMarks/ABS/internal/server/scraper/parser/utils"
	"github.com/DedAzaMarks/ABS/internal/server/statemachine"
	"github.com/DedAzaMarks/ABS/internal/server/storage"
	"log"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

const (
	ping = "ping"

	search = "new search"

	start          = "start"
	open           = "open"
	close_         = "close"
	abort          = "abort"
	help           = "help"
	ID             = "Ключ"
	NewSearch      = "Новый поиск"
	Cancel         = "Отмена"
	CancelDownload = "Отмена текущей загрузки"
)

const (
	SearchEndpoint = "https://3b5a02883www.lafa.site/torrentz/search/"
	HomeEndpoint   = "https://3b5a02883www.lafa.site"
)

var mainKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(ID),
		tgbotapi.NewKeyboardButton(NewSearch),
		tgbotapi.NewKeyboardButton(Cancel),
		tgbotapi.NewKeyboardButton(CancelDownload),
	),
)

type Server struct {
	bot   *tgbotapi.BotAPI
	redis *redis.Client
	repo  storage.Storage
}

func NewServer(repo storage.Storage, bot *tgbotapi.BotAPI, redis *redis.Client) *Server {
	return &Server{
		repo:  repo,
		redis: redis,
		bot:   bot,
	}
}

func (s *Server) Start() {
	ctx := context.Background()
	go s.registerNewClients(ctx)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := s.bot.GetUpdatesChan(u)
	for update := range updates {
		log.Print("recv new message")
		if update.Message != nil {
			message := update.Message
			userID := message.From.ID
			if message.IsCommand() {
				log.Print("it is command")
				switch message.Command() {
				case start:
					log.Print("start from", message.From.ID)
					userDTO, err := s.repo.LoadUser(ctx, userID)
					if err != nil {
						var msg tgbotapi.MessageConfig
						if !errors.Is(err, myerrors.ErrorUserNotFound) {
							log.Print(err)
							msg = tgbotapi.NewMessage(
								userID,
								"Ошибка при попытке добавить нового пользователя.")
							if _, err := s.bot.Send(msg); err != nil {
								log.Print(err)
							}
						}
					}
					if userDTO != nil {
						msg := tgbotapi.NewMessage(
							message.Chat.ID,
							fmt.Sprintf("Вы уже зарегистрированы. Ваш ключ: %q. Используйте его на устройстве, куда будут скачиваться фильмы.", userDTO.SessionKey))
						if _, err := s.bot.Send(msg); err != nil {
							log.Print(err)
						}
						continue
					}
					user := domain.NewTGUser(userID)

					if err := s.redis.Set(ctx, user.SessionKey, strconv.Itoa(int(user.UserID)), 0).Err(); err != nil {
						log.Printf("failed to save session to redis: %s", err)
						continue
					}
					if err := s.repo.SaveUser(ctx, userID, user); err != nil {
						log.Print(err)
						var msg tgbotapi.MessageConfig
						msg = tgbotapi.NewMessage(
							userID,
							"Ошибка при попытке добавить нового пользователя.")
						if _, err := s.bot.Send(msg); err != nil {
							log.Print(err)
						}
					}
					msg := tgbotapi.NewMessage(
						message.Chat.ID,
						fmt.Sprintf("Добавлен новый пользователь. Ваш Ключ: %q. Используйте его на устройстве, куда будут скачиваться фильмы.", user.SessionKey))
					if _, err := s.bot.Send(msg); err != nil {
						log.Print(err)
					}
					msg = tgbotapi.NewMessage(userID, fmt.Sprintf(
						`Для начала поиска нажмите на кнопку %q
Для того, чтобы узнать свой ID  нажмите на кнопку %q. Используйте его на устройстве, куда будут скачиваться фильмы.
Для отмены поиска/выбора фильма нажмите кнопку %q
`, NewSearch, ID, Cancel))
					msg.ReplyMarkup = mainKeyboard
					if _, err := s.bot.Send(msg); err != nil {
						log.Print(err)
					}
				}
				continue
			}
			user, err := s.repo.LoadUser(ctx, userID)
			if err != nil {
				if errors.Is(err, myerrors.ErrorUserNotFound) {
					msg := tgbotapi.NewMessage(userID, "Незарегистрированный пользователь. Для регистрации воспользуйтесь командой /start")
					if _, err := s.bot.Send(msg); err != nil {
						log.Print(err)
					}
					continue
				}
				log.Print(err)
				continue
			}
			switch message.Text {
			case ID:
				msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Ваш Ключ: %q. Используйте его на устройстве, куда будут скачиваться фильмы.", user.SessionKey))
				msg.ReplyToMessageID = message.MessageID
				if _, err := s.bot.Send(msg); err != nil {
					log.Print(err)
				}
				continue
			case NewSearch:
				if err := user.State.TriggerEvent(statemachine.EventNewSearch); err != nil {
					log.Printf("user %s error: %s", userID, err)
					user.State.Reset()
					user.SearchResults = user.SearchResults[:0]
					user.FilmResults = user.FilmResults[:0]
					if err := s.repo.SaveUser(ctx, userID, user); err != nil {
						log.Print(err)
					}
					msg := tgbotapi.NewMessage(userID, "Простите, что-то пошло не так, пожалуйста начните новый поиск.")
					if _, err := s.bot.Send(msg); err != nil {
						log.Print(err)
					}
				}
				if err := s.repo.SaveUser(context.Background(), userID, user); err != nil {
					log.Print(err)
				}
				msg := tgbotapi.NewMessage(userID, "Введите название фильма.")
				if _, err := s.bot.Send(msg); err != nil {
					log.Print(err)
				}
				continue
			case Cancel:
				user.State.Reset()
				user.SearchResults = user.SearchResults[:0]
				user.FilmResults = user.FilmResults[:0]
				if err := s.repo.SaveUser(ctx, userID, user); err != nil {
					log.Print(err)
				}
				msg := tgbotapi.NewMessage(userID, fmt.Sprintf("Для начала поиска нажмите на кнопку %q", NewSearch))
				if _, err := s.bot.Send(msg); err != nil {
					log.Print(err)
				}
				continue
			case CancelDownload:
				msg := tgbotapi.NewMessage(userID, "Если на ваше устройство идет скачивание, оно будет прервано")
				if _, err := s.bot.Send(msg); err != nil {
					log.Print(err)
				}
				cancelDownloadMessage := pb.ServerToClientChannelMessage{
					Action: &pb.ServerToClientChannelMessage_Stop{
						Stop: &pb.ServerToClientChannelMessage_StopDownload{DownloadID: ""}}}
				pbBuf, _ := proto.Marshal(&cancelDownloadMessage)
				log.Printf("publish to %s: link", user.SessionKey)
				if cmd := s.redis.Publish(ctx, user.SessionKey, pbBuf); cmd.Err() != nil {
					log.Print(cmd.Err())
				}
				continue
			}
			if user.State.CurrentState() == statemachine.StateSearch {
				title := message.Text
				go s.SearchFilm(message, user, title)
				continue
			}
			msg := tgbotapi.NewMessage(userID, "Сейчас не ожидается никакой ввод.")
			msg.ReplyToMessageID = message.MessageID
			if _, err := s.bot.Send(msg); err != nil {
				log.Print(err)
			}
			continue
		} else if update.CallbackQuery != nil {
			userID_ := update.CallbackQuery.From.ID
			userID := strconv.FormatInt(update.CallbackQuery.From.ID, 10)

			callback := tgbotapi.NewCallback(userID, update.CallbackQuery.Data)
			if _, err := s.bot.Request(callback); err != nil {
				log.Print(err)
			}
			user, err := s.repo.LoadUser(ctx, userID_)
			if err != nil {
				if errors.Is(err, myerrors.ErrorUserNotFound) {
					log.Print("user not found")
					msg := tgbotapi.NewMessage(userID_, "Незарегистрированный пользователь. Для регистрации воспользуйтесь командой /start")
					if _, err := s.bot.Send(msg); err != nil {
						log.Print(err)
					}
					continue
				}
				log.Print(err)
				continue
			}
			if user.State.CurrentState() == statemachine.StateFilmSelection {
				go s.GetFilmLinks(update.CallbackQuery, user)
				continue
			} else if user.State.CurrentState() == statemachine.StateVersionSelection {
				if len(user.Devices) == 0 {
					msg := tgbotapi.NewMessage(userID_, fmt.Sprintf("У вас не привязано устройство. Пожалуйста воспользуйтесь %q для привязки устройства, на которое будет загружен фильм.", user.SessionKey))
					if _, err := s.bot.Send(msg); err != nil {
						log.Print(err)
					}
					continue
				}
				go func() {
					magnetLink, ok := s.SelectFilmVersion(update.CallbackQuery, user)
					if !ok {
						log.Printf("user %s is not a film link", user.SessionKey)
						return
					}
					message := pb.ServerToClientChannelMessage{
						Action: &pb.ServerToClientChannelMessage_Start{
							Start: &pb.ServerToClientChannelMessage_StartDownload{
								Href: magnetLink}}}
					pbBuf, _ := proto.Marshal(&message)
					log.Printf("publish to %s: link", userID)
					if cmd := s.redis.Publish(ctx, user.SessionKey, pbBuf); cmd.Err() != nil {
						log.Print(cmd.Err())
					}
				}()
				continue
			} else {
				log.Print("call back on wrong state")
				msg := tgbotapi.NewMessage(userID_, fmt.Sprintf("Для поиска фильма нажмите %q", NewSearch))
				if _, err := s.bot.Send(msg); err != nil {
					log.Print(err)
				}
				continue
			}
		}
	}
}

func (s *Server) registerNewClients(ctx context.Context) {
	registerSub := s.redis.Subscribe(ctx, "register_new_client")
	log.Print("subscribed to register_new_client")
	defer func(registerSub *redis.PubSub) { _ = registerSub.Close() }(registerSub)
	for {
		register, err := registerSub.ReceiveMessage(ctx)
		if err != nil {
			log.Print(err)
			continue
		}
		log.Print("got register message")
		var registerReq pb.RegisterNewClient
		if err := proto.Unmarshal([]byte(register.Payload), &registerReq); err != nil {
			log.Print(err)
			continue
		}
		log.Print("register message unmarshalled")
		userID_, err := s.redis.Get(ctx, registerReq.GetSessionKey()).Result()
		if err != nil {
			log.Print(err)
			continue
		}
		userID, _ := strconv.ParseInt(userID_, 10, 64)
		clientID := uuid.MustParse(registerReq.DeviceID)
		if err := s.repo.AddNewDevice(ctx, userID, clientID, registerReq.DeviceName); err != nil {
			if errors.Is(err, myerrors.ErrorUserNotFound) {
				msg := tgbotapi.NewMessage(userID, "Незарегистрированный пользователь. Для регистрации воспользуйтесь командой /start")
				if _, err := s.bot.Send(msg); err != nil {
					log.Print(err)
				}
				continue
			}
			if _, err := s.bot.Send(tgbotapi.NewMessage(userID, "Не получилось добавить устройство")); err != nil {
				log.Print(err)
			}
			log.Print(err)
			continue
		}
		if _, err := s.bot.Send(tgbotapi.NewMessage(userID, fmt.Sprintf("Добавлено новое устройство %s", registerReq.DeviceName))); err != nil {
			log.Print(err)
		}
		log.Print("new client registered")
	}
}

func (s *Server) SearchFilm(message *tgbotapi.Message, user *domain.User, title string) {
	userID := message.From.ID
	user.SearchResults = user.SearchResults[:0]
	win, err := utils.UTF2WIN(title)
	URLEncoded := url.QueryEscape(win)
	replacedSpaces := strings.ReplaceAll(URLEncoded, "+", "%20")
	searchURL := SearchEndpoint + replacedSpaces + "/"
	searchResults, err := scraper.Search(searchURL)
	if err != nil {
		log.Printf("scraper search returned error: %s", err)
		msg := tgbotapi.NewMessage(userID, "error on title search")
		msg.ReplyToMessageID = message.MessageID
		if _, err := s.bot.Send(msg); err != nil {
			log.Print(err)
		}
		return
	}
	log.Printf("search result ok: %v", searchResults)
	if len(searchResults) == 0 {
		msg := tgbotapi.NewMessage(userID, "По такому запросу ничего не найдено. Попробуйте еще.")
		msg.ReplyToMessageID = message.MessageID
		if _, err := s.bot.Send(msg); err != nil {
			log.Print(err)
		}
		return
	}
	msg := tgbotapi.NewMessage(userID, "Результаты поиска")
	msg.ReplyToMessageID = message.MessageID
	var keyboard [][]tgbotapi.InlineKeyboardButton
	for _, searchResult := range searchResults {
		searchResultID := uuid.New()
		user.SearchResults = append(user.SearchResults, domain.SignedSearchResult{
			ID:           searchResultID,
			SearchResult: searchResult,
		})
		keyboard = append(keyboard,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(searchResult.Title, searchResultID.String())))
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)
	if _, err := s.bot.Send(msg); err != nil {
		log.Print(err)
	}
	if err := s.repo.SaveUser(context.Background(), userID, user); err != nil {
		log.Print(err)
	}
	if err := user.State.TriggerEvent(statemachine.EventSelectFilm); err != nil {
		log.Printf("user %d error: %s", userID, err)
		user.State.Reset()
		user.SearchResults = user.SearchResults[:0]
		user.FilmResults = user.FilmResults[:0]
		if err := s.repo.SaveUser(context.Background(), userID, user); err != nil {
			log.Print(err)
		}
		msg := tgbotapi.NewMessage(userID, "Простите, что-то пошло не так, пожалуйста начните новый поиск.")
		if _, err := s.bot.Send(msg); err != nil {
			log.Print(err)
		}
	}
	if err := s.repo.SaveUser(context.Background(), userID, user); err != nil {
		log.Print(err)
	}
}

func (s *Server) GetFilmLinks(callbackQuery *tgbotapi.CallbackQuery, user *domain.User) {
	userID := callbackQuery.From.ID
	searchResultID := uuid.MustParse(callbackQuery.Data)
	indx := slices.IndexFunc(user.SearchResults, func(sr domain.SignedSearchResult) bool {
		return sr.ID == searchResultID
	})
	if indx == -1 {
		log.Printf("film with id %s not found in users(%d) search results: %v", searchResultID, userID, user.SearchResults)
		if _, err := s.bot.Send(tgbotapi.NewMessage(userID, fmt.Sprintf(
			"Выбран неизвестный фильм. Выберите фильм из текущей подборки. Если же вы хотите выбрать другой фильм нажмите %q",
			NewSearch,
		))); err != nil {
			log.Print(err)
		}
		return
	}
	requestUrl := HomeEndpoint + user.SearchResults[indx].SearchResult.Href
	filmResults, err := scraper.Film(requestUrl)
	if err != nil {
		log.Print(err)
		if _, err := s.bot.Send(tgbotapi.NewMessage(userID, "Не получилось найти фильм с таким названием. Возможно это сериал или игра. Сейчас для скачивания доступны только фильмы.")); err != nil {
			log.Print(err)
		}
		return
	}
	for _, filmResult := range filmResults {
		filmResultID := uuid.New()
		user.FilmResults = append(user.FilmResults, domain.SignedFilmResult{
			ID:         filmResultID,
			FilmResult: filmResult,
		})
	}
	//var tableSize int
	var outputTable strings.Builder
	for i, result := range filmResults {
		//tmpTableSize := 0
		i++
		n := strconv.Itoa(i)
		outputTable.WriteString(n)
		outputTable.WriteString(") ")
		outputTable.WriteString(result.Quality)
		outputTable.WriteByte(' ')
		outputTable.WriteString(result.TranslationVoiceover)
		outputTable.WriteByte(' ')
		outputTable.WriteString(result.Author)
		outputTable.WriteByte(' ')
		outputTable.WriteString(result.FileFormat)
		outputTable.WriteByte(' ')
		outputTable.WriteString(result.Size)
		outputTable.WriteByte('\n')
		//if tableSize+tmpTableSize > 4096 {
		//}
	}
	msg := tgbotapi.NewMessage(userID, outputTable.String())
	msg.ReplyToMessageID = callbackQuery.Message.MessageID
	if _, err := s.bot.Send(msg); err != nil {
		log.Print(err)
	}
	msg = tgbotapi.NewMessage(userID, "Выберите версию фильма")
	var keyboard [][]tgbotapi.InlineKeyboardButton
	row := tgbotapi.NewInlineKeyboardRow()
	for i, filmResult := range filmResults {
		i++
		filmResultID := uuid.New()
		user.FilmResults = append(user.FilmResults, domain.SignedFilmResult{
			ID:         filmResultID,
			FilmResult: filmResult,
		})
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(i), filmResultID.String()))
		if i > 0 && i%5 == 0 {
			keyboard = append(keyboard, row)
			row = tgbotapi.NewInlineKeyboardRow()
		}
	}
	if len(filmResults)%5 != 0 {
		keyboard = append(keyboard, row)
	}
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(keyboard...)
	if _, err := s.bot.Send(msg); err != nil {
		log.Print(err)
	}
	if err := s.repo.SaveUser(context.Background(), userID, user); err != nil {
		log.Print(err)
	}
	if err := user.State.TriggerEvent(statemachine.EventSelectVersion); err != nil {
		log.Printf("user %d error: %s", userID, err)
		user.State.Reset()
		user.SearchResults = user.SearchResults[:0]
		user.FilmResults = user.FilmResults[:0]
		if err := s.repo.SaveUser(context.Background(), userID, user); err != nil {
			log.Print(err)
		}
		msg := tgbotapi.NewMessage(userID, "Простите, что-то пошло не так, пожалуйста начните новый поиск.")
		if _, err := s.bot.Send(msg); err != nil {
			log.Print(err)
		}
	}
	if err := s.repo.SaveUser(context.Background(), userID, user); err != nil {
		log.Print(err)
	}
}

func (s *Server) SelectFilmVersion(query *tgbotapi.CallbackQuery, user *domain.User) (string, bool) {
	userID := query.From.ID
	versionID := uuid.MustParse(query.Data)
	indx := slices.IndexFunc(user.FilmResults, func(result domain.SignedFilmResult) bool {
		return result.ID == versionID
	})
	if indx == -1 {
		log.Printf("version with id %s not found in users(%d) versions: %v", versionID, userID, user.FilmResults)
		if _, err := s.bot.Send(tgbotapi.NewMessage(userID, fmt.Sprintf(
			"Неизвестная версия фильма. Выберите версию из текущей подборки. Если же вы хотите выбрать другой фильм нажмите %q",
			NewSearch,
		))); err != nil {
			log.Print(err)
		}
		return "", false
	}
	filmResult := user.FilmResults[indx].FilmResult
	if err := user.State.TriggerEvent(statemachine.EventFinish); err != nil {
		log.Printf("user %d error: %s", userID, err)
		user.State.Reset()
		user.SearchResults = user.SearchResults[:0]
		user.FilmResults = user.FilmResults[:0]
		msg := tgbotapi.NewMessage(userID, "Простите, что-то пошло не так, пожалуйста начните новый поиск.")
		if _, err := s.bot.Send(msg); err != nil {
			log.Print(err)
		}
	}
	if err := s.repo.SaveUser(context.Background(), userID, user); err != nil {
		log.Print(err)
	}
	if _, err := s.bot.Send(tgbotapi.NewMessage(userID, fmt.Sprintf("Фильм с выбранными характеристиками (%s %s %s %s %s) будет скачан на устройство, которое вы привязали. Чтобы скачать новый фильм или выбрать другую версию нажмите на кнопку %q",
		filmResult.Quality,
		filmResult.TranslationVoiceover,
		filmResult.Author,
		filmResult.FileFormat,
		filmResult.Size,
		NewSearch,
	))); err != nil {
		log.Print(err)
	}
	return filmResult.Magnet, true
}
