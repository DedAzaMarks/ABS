document.addEventListener("DOMContentLoaded", function(event) {
	$('.rate_widget').each(function(i) {
		var widget = this;
  		set_votes(widget);
  	});
  	$('.rate_widget2').each(function(i) {
		var widget = this;
  		set_votes(widget);
  	});
	$('.ratings_stars').hover(
		function() {
			if($(this).parent().find('#already').html() == 0)
        	{
				$(this).prevAll().andSelf().addClass('ratings_over');
   				$(this).nextAll().removeClass('ratings_vote');
   			}
     	},
	    function() {
	    	if($(this).parent().find('#already').html() == 0)
	    	{
	     		$(this).prevAll().andSelf().removeClass('ratings_over');
	         	set_votes($(this).parent());
	      	}
	    }
    );
	$('.ratings_stars').bind('click', function() {
 		var star = this;
		var widget = $(this).parent();
		var currentv =widget.find('#current').html();

     	var curclick = $(star).attr('class').split(' ')[0];

        var r_tid = widget.attr('class').split(' ')[1];
        console.log(r_tid);
		widget.find('.'+curclick).prevAll().andSelf().addClass('ratings_vote');
     	widget.find('.'+curclick).nextAll().removeClass('ratings_vote');

     	widget.find('#already').html('1');
     		$.ajax({
			  type: "GET",
			  url: "/ajax.php",
			  data: { update_kp_rating: "yes", torrentid:r_tid.replace("tid_",""), user_rating: curclick.replace("star_","") }
			})
  	});

});

function set_votes(widget) {
	var avg = $(widget).find(">:first-child").html();

	$(widget).find('.star_' + avg).prevAll().andSelf().addClass('ratings_vote');
	$(widget).find('.star_' + avg).nextAll().removeClass('ratings_vote');
}