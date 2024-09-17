// ==UserScript==
// @name         aqyComment
// @namespace    http://tampermonkey.net/
// @version      2024-09-16
// @description  try to take over the world!
// @author       You
// @match        http://vedio.kundihulan.com/play/93508-1-4.html
// @icon         https://www.google.com/s2/favicons?sz=64&domain=kundihulan.com
// @grant        none
// ==/UserScript==

(function() {
	var css = `.floating {
		position: fixed;
		padding: 15px;
		animation: float 18s;
		//animation-iteration-count: 1;
		top: 0;
		right: 0;
		width: 100% ;
		test-align: right;
		z-index: 999999999;
		float: left;
		color: #fff;
		font-size: 14px;
	}
    @keyframes float {
		0%{
			left: 100% ;
		}
		100%{
			left: -50% ;
		}
	}
    #sub_form_aqy {
		position: fixed;
		padding: 15px;
		text-align: center;
		width: 100% ;
		height: 60px;
		background-color: #000;
		color: #fff;
		z-index: 9999;
		top: 0;
	}
    #sub_form_aqy * {
		padding: 5px;
		margin: 8px;
	}

	`;
	const style = document.createElement('style');
    style.appendChild(document.createTextNode(css));
	document.head.appendChild(style)
	// aiqiyi链接
	var url = '';
	var line = 0;
	window.maxLine = 7;
	//获取弹幕，省事，一次性返回
	window.contents = {};

	var div = document.createElement("div");

	div.innerHTML = "<div id='myTopTips'></div>";
	document.body.appendChild(div);

	var video;

	// 获取播放器
	var iframes = document.getElementsByTagName("iframe");

	iframes[iframes.length - 1].onload = function() {
		var divForm = document.createElement("div");
		divForm.innerHTML = "<div id='sub_form_aqy'> 弹幕行数<input type='text' id='aiqiyi_line'value=8 > <input type='text' id='aiqiyi_v' value placeholder='爱奇艺ID'><a id='get_db_btn' href='javascript:window.getDm()' >下载弹幕</a></div>";
		document.body.appendChild(divForm);
	};

	function generateUniqueID() {
		const timestamp = Date.now().toString(36);
		const uniqueRandom = Math.random().toString(36).substr(2, 9);
		return `${timestamp}-${uniqueRandom}`;;
	}

	function getTipsDiv() {
		// return document.body
		if (window.fullScreen) {
			return video.parentNode;
		} else {
			return document.getElementById('myTopTips');
		}
	}

	function writeDiv(content) {
		var div = document.createElement("div");
		if (line >= window.maxLine) {
			line = 0;
		}
		div.style = "top:" + (line * 25) + "px;";
		let uniqueID = generateUniqueID();
		div.classList.add("floating");
		line = line + 1;

		div.id = uniqueID;
		div.innerHTML = content;

		getTipsDiv().appendChild(div);
		div.addEventListener('animationend',
		function() {
			this.remove();
		});
	}
	function start() {
		for (var i = iframes.length - 1; i >= 0; i--) {
			// 获取iframe的document对象
			var iframeDoc = iframes[i].contentDocument || iframes[i].contentWindow.document;
			// 获取video标签
			var v = iframeDoc.getElementsByTagName("video");
			//alert("视频数量" + v.length)
			// 你现在可以操作这个video标签了，例如播放视频
			if (v.length > 0) {
				//iframe 添加css
				const style = document.createElement('style');
                style.appendChild(document.createTextNode(css));
				iframeDoc.head.appendChild(style);
				// 绑定iframe
				window.videoIframe = iframeDoc;
				video = v[0];
				//alert("视频已加载完成");
				writeDiv("找到播放器，可以开始播放");
				//video.play()
				//writeDiv("视频已加载完成!");
				done();
                video.addEventListener('loadedmetadata',
				function() {
					//alert("视频已加载完成");
					//writeDiv("视频已加载完成!");
					//done()
				});
			}
		}
	}

	async function putContent(contents, st) {
		var count = contents.length;
		var sleepTime = (st / count).toFixed(0)
		console.log(sleepTime);
		//getTipsDiv().appendChild(div);
		for (var i = 0; i < contents.length; i++) {
			//alert(contents[i])
			writeDiv(contents[i])
			await sleep(sleepTime);
		}
	}

	//避免同一个key 重复弹出
	var lastTIme = 0;
	async function done() {
		//根据播放速度调整时间
		var st = (1000 / video.playbackRate).toFixed(0);
		if (!video.paused) {
			console.log('视频正在播放');
			//alert("视频正在播放")
			var time = (video.currentTime).toFixed(0);
			var k = "_" + time
			if (window.contents.hasOwnProperty(k) && time != lastTIme) {
				lastTIme = time;
				var currentContents = window.contents[k];
				console.log(currentContents);
				putContent(currentContents, st);
			}
		}
		setTimeout(done, st)
	}

	function sleep(ms) {
		return new Promise(resolve =>setTimeout(resolve, ms));
	}
	window.dmLoadIng = false;
    window.getDm = function() {
		if (window.dmLoadIng) {
			return
		}
		let v = document.getElementById("aiqiyi_v").value;
		let l = document.getElementById("aiqiyi_line").value;
		if (!v) {
			alert("aiqiyi_v 不能为空");
            return false
		}
		window.dmLoadIng = true;
		document.getElementById("get_db_btn").text = "下载中";
		let xhr = new XMLHttpRequest();
		xhr.open("GET", "http://127.0.0.1:1188/?id=" + v, true);
		xhr.onreadystatechange = function() {
			if (xhr.readyState == 4 && xhr.status == 200) {
				var response = JSON.parse(xhr.responseText);
				window.maxLine = l;
				window.contents = response;
				//alert(window.contents)
				document.getElementById('sub_form_aqy').style = "display:none";
				start()

			}
		};
		xhr.send();
	}

})();

function checkFullScreen() {
	var isFull = document.fullscreenElement || document.mozFullScreenElement || document.webkitFullscreenElement || document.msFullscreenElement;
	return isFull === null || isFull === undefined ? false: true;
}
window.fullScreen = false;
window.videoIframe = false;
window.onresize = function() {
	if (checkFullScreen()) {
		window.fullScreen = true;
	} else {
		window.fullScreen = false;
		if (window.videoIframe != false) {
			// TODO 退出全屏，里面的弹幕删除不了
			var elements = window.videoIframe.getElementsByClassName("floating");
            for (var i = 0; i < elements.length; i++) {
				//var element = elements[i];
				//window.videoIframe.getElementById(element.id).remove()
				//if (element){
				//element.parentNode.removeChild(element);
				//}
				//alert(element.id)
				//element.parentNode.removeChild(element);
			}
		}
	}
}