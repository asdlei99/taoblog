/* 回到顶端 */
window.onscroll = function() {
	if(window.scrollY > 160) {
		$("#back-to-top").fadeIn(500);
	} else {
		$("#back-to-top").fadeOut(500);
	}
};
$('#back-to-top').click(function(){
	$('html,body').animate({
		scrollTop: 0
	}, 300);
});

/*  阅读模式 */
var in_reading_mode = false;

function toggle_reading_mode() {
	var header = $('#header');
	var main = $('#main');
	var icon = $('#reading-mode i');

	if(!in_reading_mode) {
		header.css('left', '-300px');
		main.css('margin-left', '0px');
		icon.css('color', '#f66');
	} else {
		header.css('left', '0px');
		main.css('margin-left', '300px');
		icon.css('color', 'inherit');
	}

	in_reading_mode = !in_reading_mode;

	show_tips(in_reading_mode 
		? '<b>已进入阅读模式。</b><br/>您可以点击右下角的“+”号退出阅读模式。'
		: '已退出阅读模式。');
}

$('#reading-mode').click(function() {
	toggle_reading_mode();
});

/* 字体大小调整 */
$('.font-sizing').click(function(e){
	var post = $('.entry');
	var cl = e.target.classList;
	if(cl.contains('inc') || cl.contains('fa-plus')) {
		var newSize = parseFloat(post.css('font-size')) * 1.2 + 'px';
		post.css('font-size', newSize);
		if(window.localStorage) localStorage.setItem('font-size', newSize);
		if(typeof show_tips == 'function') show_tips('字体大小: '+newSize);
	} else if(cl.contains('dec') || cl.contains('fa-minus')) {
		var newSize = Math.max(8, parseFloat(post.css('font-size')) / 1.2) + 'px';
		post.css('font-size', newSize);
		if(window.localStorage) localStorage.setItem('font-size', newSize);
		if(typeof show_tips == 'function') show_tips('字体大小: '+newSize);
	}
});

/* 主页链接 */
$('.home-a').click(function() {
	location.href = location.protocol + '//' + location.host;
});

/* RSS订阅提示 */
(function() {
	if(!window.localStorage) return;

	if(document.referrer == '') {
		var ert = parseInt(localStorage.getItem('empty_referrer_times') || 0);
		if( ert == -1) return;

		ert += 1;
		localStorage.setItem('empty_referrer_times', ert);

		if(ert >= 15) {
			show_tips({
				timeout: 5000,
				content: '<b>RSS订阅</b><br/>亲，你已经直接访问本博客 <b>' + ert + '</b> 次啦，考虑<a href="/rss" target="_blank" style="color: blue;">订阅</a>本博客吧～(<a class="dont" style="cursor: pointer; color: blue;">不再提示</a>)',
				click: function(ob) {
					var target = ob.event.target;
					if(target.classList.contains('dont')) {
						localStorage.setItem('empty_referrer_times', -1);
						ob.dismiss();
					}
				}}
			);
		}
	}
})();

/* 点击图片放大 & 拖动浏览*/
/* 写得超级烂，完全没管性能 */
(function() {
	var body = $('body');
	var imgdiv = $('#img-view');

	function view_image(ele, show) {
        var img = $('#img-view img');
		if(show) {
			body.css('max-height', window.innerHeight);
			body.css('overflow', 'hidden');
			img.attr('src', ele.src);
            img.css('left', (parseInt(imgdiv.css('width'))-parseInt(img.prop('naturalWidth')))/2 + 'px');
            img.css('top', (parseInt(imgdiv.css('height'))-parseInt(img.prop('naturalHeight')))/2 + 'px');
			imgdiv.show();
		} else {
            // 以下两行清除因拖动导致的设置
			img.css('left', '0px');
			img.css('top', '0px');
			body.css('max-height', 'none');
			body.css('overflow', 'auto');

            imgview.dragging = false;
			imgdiv.hide();
		}
	}

	$('.entry img').click(function(e) {
		view_image(this, true);

		var f = function(e) {
			if(e.keyCode == 27 || e.keyCode == 32) {
				view_image(null, false);
				window.removeEventListener('keydown', f);
                e.preventDefault();
			}
		};

		window.addEventListener('keydown', f);
	});

	imgdiv.click(function() {
		view_image(null, false);
	});

    window.imgview = {};
    imgview.dragging = false;

    $('#img-view img').on('mousedown', function(e) {
        var target = e.target;
        imgview.offset_x = e.clientX;
        imgview.offset_y = e.clientY;

        imgview.coord_x = parseInt(target.style.left);
        imgview.coord_y = parseInt(target.style.top);

        imgview.dragging =true;

        e.preventDefault();
        return false;
    });

    $('#img-view img').on('mousemove', function(e) {
        if(!imgview.dragging) return;

        var target = e.target;
        target.style.left = imgview.coord_x + e.clientX - imgview.offset_x + 'px';
        target.style.top = imgview.coord_y + e.clientY - imgview.offset_y + 'px';

        e.preventDefault();
        return false;
    });

    $('#img-view img').on('mouseup', function(e) {

        imgview.dragging = false;
        e.preventDefault();
        return false;
    });

    $('#img-view img').on('click', function(e) {
        e.preventDefault();
        return false;
    });

    $('#img-view img').on('dblclick', function(e) {
        view_image(null, false);
        
        e.preventDefault();
        return false;
    });
})();

/* 目录展开与隐藏 */
(function() {
	if($('.entry .toc').length == 0) return;

	function hide_toc(hide) {
		var hide_toc = $('#hide-toc');
		var toc_ul = $('.entry .toc > ul');

		if(hide) {
			hide_toc.text('[显示]');
			toc_ul.hide();
		} else {
			hide_toc.text('[隐藏]');
			toc_ul.show();
		}
	}

	$('.entry .toc h2, .entry .toc h3').replaceWith('<div style="margin-bottom: -10px;"><h2 style="float: left; margin-right: 20px;">目录</h2><span id="hide-toc" class="no-sel" style="float: right; cursor: pointer;"></span><div style="clear: both;"></div></div>');

	hide_toc(true);

	$('#hide-toc').click(function() {
		var hidden = $('.entry .toc > ul').css('display') == 'none';
		hide_toc(!hidden);
	});

    window.addEventListener('keyup', function(e) {
        if(e.keyCode == 27) {
            hide_toc(true);
        }
    });
})();

/* pre 的双击全选与全窗口 */
(function() {
    $('.entry pre').on('dblclick', function(e) {
        if(e.target.tagName == 'PRE') {
            var selection = window.getSelection();
            var range = document.createRange();
            range.selectNodeContents(e.target);
            selection.removeAllRanges();
            selection.addRange(range);
            e.preventDefault();
            return;
        }
    });
})();

/* 服务器运行时间 */
// http://www.htmlgoodies.com/html5/javascript/calculating-the-difference-between-two-dates-in-javascript.html
(function() {
	function date_between(date1, date2) {
		// get one day in milliseconds
		var one_day = 24*60*60*1000;

		// convert to milliseconds
		var date1_ms = date1.getTime();
		var date2_ms = date2.getTime();

		// calc diff
		var diff_ms = date2_ms - date1_ms;

		return Math.round(diff_ms / one_day);
	}

	var start = new Date(2014, 12-1, 24);
	var now = new Date();

	var days = date_between(start, now);

	$('#server-run-time').text(days);
})();

