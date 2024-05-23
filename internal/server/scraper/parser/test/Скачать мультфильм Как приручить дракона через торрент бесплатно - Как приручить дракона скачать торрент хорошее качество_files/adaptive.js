jQuery(document).ready(function($){
	
	$('.mob-m button').click(function(){
		$(this).toggleClass('openm');
		$('body').toggleClass('open-m');
		if($(this).hasClass('openm')){
			$('.to_top_box').after('<div class="over_bg"></div>');
			$('.over_bg').click(function(){
				$(this).remove();
				$('body').toggleClass('open-m');
				$('.mob-m button').toggleClass('openm');
			});
		}else{
			$('.over_bg').remove();
		}
	});
	
	$('.login').click(function(){
		$(this).toggleClass('openl');
		$('body').toggleClass('open-l');
		if($(this).hasClass('openl')){
			$('.to_top_box').after('<div class="over_bg"></div>');
			$('.over_bg').click(function(){
				$(this).remove();
				$('body').toggleClass('open-l');
				$('.login').toggleClass('openl');
			});
		}else{
			$('.over_bg').remove();
		}
	});
	
	var nav = $('.mob-m');
 
	$(window).scroll(function () {
		if ($(this).scrollTop() > 50) {
			nav.addClass("fix");
		} else {
			nav.removeClass("fix");
		}
	});
	tableRemove();
	$(window).resize(tableRemove);

	if($('.finger').length < 1){
		$('#tablesorter').before('<div class="finger"></div>');
	}
});

function tableRemove(){
	if($(window).outerWidth() < 640){
		$('table.cinema_box').each(function(){
			$(this).replaceWith( $(this).html()
				.replace(/<tbody/gi, "<div class='new_box_cinema'")
				.replace(/<tr>/gi, "")
				.replace(/<\/tr>/gi, "")
				.replace(/<td/gi, "<div class='s_c_box'")
				.replace(/<\/td>/gi, "</div>")
				.replace(/<\/tbody/gi, "<\/div")
			);
		});

		$('.s_c_box').find('.seprat').each(function(){
			$(this).parent().remove();
		});



		$(".expand-child td").html(function (i, html) {
			return html.replace(/&nbsp;/g, '');
		});
		if($('.finger').length < 1){
			$('#tablesorter').before('<div class="finger"></div>');
		}

	}
}