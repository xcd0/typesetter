
/*
window.addEventListener('DOMContentLoaded', () => {
var ps = document.getElementsByTagName('p');
for (var p of ps) {
	var out = "";
	let s = p.innerHTML;
	for (var c of s) {
		out += "<span>" + c + "</span>";
	};
	p.innerHTML = out;
};

})

*/


function checkDrawSize(html){
	var span = document.createElement('span');
	//span.style.top = '-1000px';
	//span.style.left = '-1000px';
	span.style.position = 'absolute';
	span.style.whiteSpace = 'nowrap';
	span.style.fontSize = '16px';
	span.style.lineHeight= '1em';
	span.style.letterSpacing = '0';
	span.innerHTML = html
	document.body.appendChild(span);
	for( var h of document.querySelectorAll('h1,h2,h3,h4,h5,h6,h7')){
		h.style.fontSize = '16px';
		h.style.lineHeight = '1em';
		h.style.letterSpacing = '0';
	}
	var width = span.clientWidth;
	var height = span.clientHeight;
	span.parentElement.removeChild(span);
	return [width, height];
}

window.addEventListener('DOMContentLoaded', () => {
	var ps = document.getElementsByClassName('page');
	console.log("ps.length:",ps.length)
	let regtagstart =/<\w+(?:\s+\w+(?:\s*=\s*(?:"[^"]*"|'[^']*'|[\^'">\s]+))?)*\s*>/;
	let regtagend = /<\/\w+\s*>/;
	let regrubystart = /<ruby>/; let regrubyend = /<\/ruby>/;
	let regheaderend = /<\/h[1-9]>/;
	var colHeight = 680 // ps[0].clientHeightから計算したい
	for (var s of ps) {
		var parsed = s.innerHTML.match(/<[^>]*>|[^<>]+/g);
		var lines = [];
		var tokens = "";
		var pretagstart = false; // 直前のトークンが開始タグだったかどうか
		var pretagend = false; // 直前のトークンが終了タグだったかどうか
		var inRuby = false; // ルビのタグないかどうか
		var size = []    // checkDrawSize()の戻り値を保持する。
		var tagqueue= [] // 開始タグを一時的に保持する
		for (var t of parsed) {
			t = t.replace(/\r?\n/g, '');
			t = t.replace(/　/g, '&nbsp;&nbsp;');

			if(t.length == 0) continue;
			let tagstart = regtagstart.test(t)
			let tagend = regtagend.test(t)
			let tagrubyend = regrubyend.test(t)
			if( regrubystart.test(t) ) { inRuby = true; }
			//if( tagstart ) {
			//	tagqueue.push(t);
			//	continue;
			//}
			if( regheaderend.test(t) ){ // </h1>などは必ず改行
				// 行が増えた
				tokens += t;
				lines.push(tokens);
				tokens = "";
				continue;
			}

			var tmp = tokens
			if( !tagstart && !tagend && !inRuby ){
				// ただの本文
				while (tagqueue.length>0) { tokens += tagqueue.pop(); } // 貯めていた開始タグを吐き出す。
				for (var i=0;i < t.length; i++) {
					var c = t[i];
					tokens = tmp;
					tmp += c;
					size_ = checkDrawSize(tokens);
					size = checkDrawSize(tmp);
					console.log('width:', size[0], ', height:', size[1], ', tmp : ', tmp);
					if ( size[0] > 26 ) { // 行が増えた ふりがななしで16 ありで26 この辺は計算式が不明
						lines.push(tokens);
						tmp = c;
					}
					if ( size[1] >= colHeight ) { // size[1]文字数
						// 行が増えた
						if( c == "。" || c == "、" || c == "」" || c == "）" ) {
							lines.push(tokens+c);
							tmp = "";
						}else{
							lines.push(tokens);
							tmp = c;
						}
						//}else if ( size[0] == prewidth ) {
						//	lines[lines.length-1] += c
					}else if(size[0] == 1){
						console.log(tmp);
					}
					if( i==t.length-1){
						tokens += c;
					}
				}
			}else{
				while (tagqueue.length>0) { tokens += tagqueue.pop(); } // 貯めていた開始タグを吐き出す。
				tokens += t;
				// 1token毎に足していって改行されたところで分割する。
				size = checkDrawSize(tokens);
				console.log('width:', size[0], ', height:', size[1], ', tokens : ', tokens);
				if(!inRuby && tagend){ // rubyの中でない終了タグ
					if ( size[0] > 16 ) { // 行が増えた
						lines.push(tokens);
						tokens = c;
					}
					if ( size[1] >= colHeight ) { // size[1]文字数
						// 行が増えた
						lines.push(tokens);
						tokens = c;
					}
				}
			}

			if( regrubyend.test(t) ) { inRuby = false; }
			pretagstart = tagstart
			pretagend = tagend
		}
	}
	//const swiper = new Swiper(".swiper", {});
});





