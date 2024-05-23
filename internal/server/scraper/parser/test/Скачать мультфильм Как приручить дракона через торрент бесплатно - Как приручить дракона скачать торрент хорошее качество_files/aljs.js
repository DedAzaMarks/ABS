/*! js-cookie v3.0.5 | MIT */
!function(e,t){"object"==typeof exports&&"undefined"!=typeof module?module.exports=t():"function"==typeof define&&define.amd?define(t):(e="undefined"!=typeof globalThis?globalThis:e||self,function(){var n=e.Cookies,o=e.Cookies=t();o.noConflict=function(){return e.Cookies=n,o}}())}(this,(function(){"use strict";function e(e){for(var t=1;t<arguments.length;t++){var n=arguments[t];for(var o in n)e[o]=n[o]}return e}var t=function t(n,o){function r(t,r,i){if("undefined"!=typeof document){"number"==typeof(i=e({},o,i)).expires&&(i.expires=new Date(Date.now()+864e5*i.expires)),i.expires&&(i.expires=i.expires.toUTCString()),t=encodeURIComponent(t).replace(/%(2[346B]|5E|60|7C)/g,decodeURIComponent).replace(/[()]/g,escape);var c="";for(var u in i)i[u]&&(c+="; "+u,!0!==i[u]&&(c+="="+i[u].split(";")[0]));return document.cookie=t+"="+n.write(r,t)+c}}return Object.create({set:r,get:function(e){if("undefined"!=typeof document&&(!arguments.length||e)){for(var t=document.cookie?document.cookie.split("; "):[],o={},r=0;r<t.length;r++){var i=t[r].split("="),c=i.slice(1).join("=");try{var u=decodeURIComponent(i[0]);if(o[u]=n.read(c,u),e===u)break}catch(e){}}return e?o[e]:o}},remove:function(t,n){r(t,"",e({},n,{expires:-1}))},withAttributes:function(n){return t(this.converter,e({},this.attributes,n))},withConverter:function(n){return t(e({},this.converter,n),this.attributes)}},{attributes:{value:Object.freeze(o)},converter:{value:Object.freeze(n)}})}({read:function(e){return'"'===e[0]&&(e=e.slice(1,-1)),e.replace(/(%[\dA-F]{2})+/gi,decodeURIComponent)},write:function(e){return encodeURIComponent(e).replace(/%(2[346BF]|3[AC-F]|40|5[BDE]|60|7[BCD])/g,decodeURIComponent)}},{path:"/"});return t}));
/* spoiler */
function showspoiler(id){
    var text = document.getElementById(id);
    var pic = document.getElementById('pic' + id);

    if(text.style.display == 'none')
    {
        text.style.display = 'block';
        pic.src = '/pic/minus.gif';
        pic.title = 'Скрыть';
    }
    else
    {
        text.style.display = 'none';
        pic.src = '/pic/plus.gif';
        pic.title = 'Показать';
    }
}

function toggleDarkMode(){
    $('html').toggleClass('ublack');
    Cookies.set("has_black_theme",(Cookies.get('has_black_theme') === 'yes' ? 'no' : 'yes'),{ expires: 1000, domain: '.lafa.site' });
}

function MM_swapImgRestore() { //v3.0
    var i,x,a=document.MM_sr; for(i=0;a&&i<a.length&&(x=a[i])&&x.oSrc;i++) x.src=x.oSrc;
}
function MM_preloadImages() { //v3.0
    var d=document; if(d.images){ if(!d.MM_p) d.MM_p=new Array();
        var i,j=d.MM_p.length,a=MM_preloadImages.arguments; for(i=0; i<a.length; i++)
            if (a[i].indexOf("#")!=0){ d.MM_p[j]=new Image; d.MM_p[j++].src=a[i];}}
}

function MM_findObj(n, d) { //v4.01
    var p,i,x;  if(!d) d=document; if((p=n.indexOf("?"))>0&&parent.frames.length) {
        d=parent.frames[n.substring(p+1)].document; n=n.substring(0,p);}
    if(!(x=d[n])&&d.all) x=d.all[n]; for (i=0;!x&&i<d.forms.length;i++) x=d.forms[i][n];
    for(i=0;!x&&d.layers&&i<d.layers.length;i++) x=MM_findObj(n,d.layers[i].document);
    if(!x && d.getElementById) x=d.getElementById(n); return x;
}

function MM_swapImage() { //v3.0
    var i,j=0,x,a=MM_swapImage.arguments; document.MM_sr=new Array; for(i=0;i<(a.length-2);i+=3)
        if ((x=MM_findObj(a[i]))!=null){document.MM_sr[j++]=x; if(!x.oSrc) x.oSrc=x.src; x.src=a[i+2];}
}


$(document).ready(function(){
    $('.flood').toggle(function(){
        $('#hflood'+$(this).attr("id")).slideDown();
    }, function(){
        $('#hflood'+$(this).attr("id")).slideUp();
    });



});



function go_more()
{
    $('#hide_more').toggle();
    $('#hide_a').hide();
}

function shareWindow(url,service)
{
    window.open(url,'Поделится','width=600,height=350');

}

jQuery(function($) { $.extend({
    form: function(url, data, method) {
        if (method == null) method = 'POST';
        if (data == null) data = {};

        var form = $('<form>').attr({
            method: method,
            action: url
        }).css({
            display: 'none'
        });

        var addData = function(name, data) {
            if ($.isArray(data)) {
                for (var i = 0; i < data.length; i++) {
                    var value = data[i];
                    addData(name + '[]', value);
                }
            } else if (typeof data === 'object') {
                for (var key in data) {
                    if (data.hasOwnProperty(key)) {
                        addData(name + '[' + key + ']', data[key]);
                    }
                }
            } else if (data != null) {
                form.append($('<input>').attr({
                    type: 'hidden',
                    name: String(name),
                    value: String(data)
                }));
            }
        };

        for (var key in data) {
            if (data.hasOwnProperty(key)) {
                addData(key, data[key]);
            }
        }

        return form.appendTo('body');
    }
}); });



function movie_seen(id)
{
    $("#movie_seen_"+id).html('<a class=\"menu_janr8\" href=\"#\" onClick=\"movie_un_seen('+id+');return false;\"><img style=\"vertical-align: bottom;\" width=\"16\" src=\"/pic/yes.png\"> <b>Уже смотрел</b></a>');
    $.ajax({
        type: "GET",
        url: "/ajax.php",
        data: { action: "set_movie_seen", has_seen: "1", tid:id }
    })
        .done(function( msg ) { });

}


function movie_un_seen(id)
{
    $("#movie_seen_"+id).html('<a class=\"menu_janr8\" href=\"#\" onClick=\"movie_seen('+id+');return false;\">Уже смотрел</a>');
    $.ajax({
        type: "GET",
        url: "/ajax.php",
        data: {  action: "set_movie_seen",  has_seen: "0", tid:id }
    })
        .done(function( msg ) { });

}

function c_up(comm_id, isLoggedIn)
{
    if(!isLoggedIn){
        return alert('Необходимо авторизироваться!');
    }
    $("#c_l_"+comm_id).removeClass('c_like_up').addClass('c_like_up_selected');
    $("#c_l_c_"+comm_id).text(parseInt($("#c_l_c_"+comm_id).text(),10)+1);
    $("#c_l_"+comm_id).prop('onclick',null).off('click');
    $.ajax({
        type: "GET",
        url: "/ajax.php",
        data: { action: "comment_like", comment_id: comm_id}
    })
        .done(function( msg ) { });
}


function c_down(comm_id, isLoggedIn)
{
    if(!isLoggedIn){
        return alert('Необходимо авторизироваться!');
    }
    $("#c_d_"+comm_id).removeClass('c_like_down').addClass('c_like_down_selected');
    $("#c_d_c_"+comm_id).text(parseInt($("#c_d_c_"+comm_id).text(),10)+1);
    $("#c_d_"+comm_id).prop('onclick',null).off('click');
    $.ajax({
        type: "GET",
        url: "/ajax.php",
        data: { action: "comment_dislike", comment_id: comm_id}
    })
        .done(function( msg ) { });
}



function movie_add_favorite(id)
{
    $("#fav_div").html("<img align=\"absmiddle\" width=\"16\" src=\"/pic/fav.png\"> <a class=\"menu_janr8\" href=\"/favorites.php\"><b>В избранном</b></a> [<a href=\"#\" onClick=\"movie_delete_favorite(" + id + "); return false;\"><font color=\"#ff0000\">Удалить</font></a>]");
    $.ajax({
        type: "GET",
        url: "/favorites.php",
        data: { ajax: "yes", action: "add", torrent_id:id }
    })
        .done(function( msg ) { });
}


function movie_delete_favorite(id)
{
    $("#fav_div").html("<img align=\"absmiddle\" width=\"16\" src=\"/pic/fav.png\"> <a onClick=\"movie_add_favorite(" + id + "); return false;\" class=\"menu_janr8\" href=\"#\">Добавить в избранное</a></div>");
    $.ajax({
        type: "GET",
        url: "/favorites.php",
        data: {  ajax: "yes", action: "delete", torrent_id:id }
    })
        .done(function( msg ) { });

}

function abp_warning_close()
{
    $(".toaster-container-partial-component").hide();
    Cookies.set("abp_warning_close","yes",{expires: 2});
    $.get( "/ajax.php", { act: "abp_close" });
}

$(function () {

    $("#y_frm").attr('src','/ya.htm');
    $(".help_top").click(function () {
        window.open("/help.php", "_blank");
    });

    $(".back_top").click(function () {

        $(window).scrollTop(0);
        return false;
    });


    $(window).scroll(function () {

        if ($(this).scrollTop() > 200) {
            $('.back_top').stop(true, true).fadeIn();
        } else {
            $('.back_top').stop(true, true).fadeOut();
        }
    });

    $('.ttable_comm').each(function() {
        $(this).tooltip({
            show: 'false', track: 'true', hide: 'false', tooltipClass: "custom-tooltip-styling-ttable",

            content: function () {
                return "<div class=\"tooltip-header\">Комментарии</div><div class=\"tooltip-separator\"></div><div class=\"tooltip-description\"><span style=\"color:#09b207;\">положительные</span> / <span style=\"color:#b3b1b1\">нейтральные</span> / <span style=\"color:red\">отрицательные</span></div>";
            }

        });
    })

});


function addFav(tid)
{		$.get( "/favorites.php?action=add&torrent_id="+tid, function( data ) {

    $('#fav_button_'+tid).html('<a class="menu_janr8" href="/favorites.php"><b>В избранном</b></a>');
});
}

$(function () {
    $("img.lazy").lazyload();
});