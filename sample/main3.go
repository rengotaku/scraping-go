package main

import (
	"bytes"
	"log"
	"net/http"

	"github.com/gocolly/colly/v2"
)

type pageInfo struct {
	StatusCode int
	Links      map[string]int
	Body       []byte
	DomHtml    string
}

func handler(w http.ResponseWriter, r *http.Request) {
	URL := r.URL.Query().Get("url")
	if URL == "" {
		log.Println("missing URL argument")
		return
	}
	log.Println("visiting", URL)

	c := colly.NewCollector()
	c.UserAgent = r.Header.Get("User-Agent")

	p := &pageInfo{Links: make(map[string]int)}

	// count links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href") + ":" + e.ChildText("."))
		if link != "" {
			p.Links[link]++
		}
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		// e.DOM.Find("script").Each(func(i int, s *goquery.Selection) {
		// 	s.Remove()
		// })
		// e.DOM.Find("body").Parent().Each(func(_ int, s *goquery.Selection) {
		// 	s.AddNodes()
		// })
		p.DomHtml, _ = e.DOM.Html()
	})

	// extract status code
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		p.StatusCode = r.StatusCode
		p.Body = r.Body
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
		p.StatusCode = r.StatusCode
	})

	c.Visit(URL)

	w.Header().Add("Content-Type", "text/html")

	// w.Write(p.Body)
	w.Write(bytes.NewBufferString(p.DomHtml).Bytes())

	// body := string(p.Body)
	// i := strings.LastIndex(body, "</body>")
	// newBody := body[:i] + addScript() + "</body></html>"
	// w.Write(bytes.NewBufferString(newBody).Bytes())

	// dump results
	// b, err := json.Marshal(p)
	// if err != nil {
	// 	log.Println("failed to serialize response:", err)
	// 	return
	// }
	// w.Header().Add("Content-Type", "application/json")
	// w.Write(b)
}

func main() {
	// example usage: curl -s 'http://127.0.0.1:7171/?url=http://go-colly.org/'
	addr := ":7171"

	http.HandleFunc("/", handler)

	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func addScript() string {
	return `
<script>
window.addEventListener('load', (event) => {
  var style = document.createElement('style');
  var head = document.getElementsByTagName('head');
  style.setAttribute('type', 'text\/css');
  style.innerHTML = ` + "`" + `
  *.coloring-depth {
      background-color:rgba(255,0,0,.2)!important;
      opacity: 0.97;
  }
  *.coloring-depth * {
      background-color:rgba(0,255,0,.2)!important;
  }
  *.coloring-depth * * {
      background-color:rgba(0,0,255,.2)!important;
  }
  *.coloring-depth * * * {
      background-color:rgba(255,0,255,.2)!important;
  }
  *.coloring-depth * * * * {
      background-color:rgba(0,255,255,.2)!important;
  }
  *.coloring-depth * * * * * {
      background-color:rgba(255,255,0,.2)!important;
  }
  *.coloring-depth * * * * * * {
      background-color:rgba(255,0,0,.2)!important;
  }
  *.coloring-depth * * * * * * * {
      background-color:rgba(0,255,0,.2)!important;
  }
  *.coloring-depth * * * * * * * * {
      background-color:rgba(0,0,255,.2)!important;
  }
  ` + "`" + `;
  style.innerHTML += ` + "`" + `
  #tooltip-7c392fb14fe25f428f3194f59b5b01e1c6adf8702e41755abb774812de3238dc {
    background-color: #333;
    color: white;
    padding: 5px 10px;
    border-radius: 4px;
    font-size: 13px;
  }
  ` + "`" + `;
  head[0].appendChild(style);

  // window.onbeforeunload = function(){
  //     return 'Are you sure you want to leave?';
  // };

  document.onclick = function (e) {
    e = e ||  window.event;
    var element = e.target || e.srcElement;

		if(!e.classList.contains("631181e4fd098a70d1f99c00192da69f9e5748935f85479e76786cc74d2772bb")) {
      return false;
    }
    return true;
  };

  /* https://github.com/fczbkk/css-selector-generator */
  (function(){var u=[].indexOf||function(t){for(var e=0,n=this.length;e<n;e++)if(e in this&&this[e]===t)return e;return-1};function t(t){null==t&&(t={}),this.options={},this.setOptions(this.default_options),this.setOptions(t)}t.prototype.default_options={selectors:["tag","id","class","nthchild"]},t.prototype.setOptions=function(t){var e,n,r;for(e in null==t&&(t={}),r=[],t)n=t[e],this.default_options.hasOwnProperty(e)?r.push(this.options[e]=n):r.push(void 0);return r},t.prototype.isElement=function(t){return!(1!==(null!=t?t.nodeType:void 0))},t.prototype.getParents=function(t){var e,n;if(n=[],this.isElement(t))for(e=t;this.isElement(e);)n.push(e),e=e.parentNode;return n},t.prototype.getTagSelector=function(t){return t.tagName.toLowerCase()},t.prototype.sanitizeItem=function(t){return escape(t).replace(/\%/g,"\\")},t.prototype.validateId=function(t){return null!=t&&!/^\d/.exec(t)&&1===document.querySelectorAll("#"+t).length},t.prototype.getIdSelector=function(t){var e;return null!=(e=t.getAttribute("id"))&&(e=this.sanitizeItem(e)),e=this.validateId(e)?"#"+e:null},t.prototype.getClassSelectors=function(t){var o,l,e;return e=[],null!=(o=t.getAttribute("class"))&&""!==(o=(o=o.replace(/\s+/g," ")).replace(/^\s|\s$/g,""))&&(e=function(){var t,e,n,r;for(r=[],t=0,e=(n=o.split(/\s+/)).length;t<e;t++)l=n[t],r.push("."+this.sanitizeItem(l));return r}.call(this)),e},t.prototype.getAttributeSelectors=function(t){var e,n,r,o,l,i,s;for(r=[],n=["id","class"],o=0,l=(i=t.attributes).length;o<l;o++)s=(e=i[o]).nodeName,u.call(n,s)<0&&r.push("["+e.nodeName+"="+e.nodeValue+"]");return r},t.prototype.getNthChildSelector=function(t){var e,n,r,o,l,i;if(null!=(n=t.parentNode))for(l=e=0,i=(o=n.childNodes).length;l<i;l++)if(r=o[l],this.isElement(r)&&(e++,r===t))return":nth-child("+e+")";return null},t.prototype.testSelector=function(t,e){var n,r;return n=!1,null!=e&&""!==e&&1===(r=t.ownerDocument.querySelectorAll(e)).length&&r[0]===t&&(n=!0),n},t.prototype.getAllSelectors=function(t){var e;return e={t:null,i:null,c:null,a:null,n:null},0<=u.call(this.options.selectors,"tag")&&(e.t=this.getTagSelector(t)),0<=u.call(this.options.selectors,"id")&&(e.i=this.getIdSelector(t)),0<=u.call(this.options.selectors,"class")&&(e.c=this.getClassSelectors(t)),0<=u.call(this.options.selectors,"attribute")&&(e.a=this.getAttributeSelector(t)),0<=u.call(this.options.selectors,"nthchild")&&(e.n=this.getNthChildSelector(t)),e},t.prototype.testUniqueness=function(t,e){var n;return 1===(n=t.parentNode.querySelectorAll(e)).length&&n[0]===t},t.prototype.getUniqueSelector=function(t){var e,n,r;if(null!=(r=this.getAllSelectors(t)).i)return r.i;if(this.testUniqueness(t,r.t))return r.t;if(0!==r.c.length){if(n=e=r.c.join(""),this.testUniqueness(t,n))return n;if(n=r.t+e,this.testUniqueness(t,n))return n}return r.n},t.prototype.getSelector=function(t){var e,n,r,o,l,i,s,u,c,p;for(e=[],s=0,c=(r=this.getParents(t)).length;s<c;s++)n=r[s],null!=(l=this.getUniqueSelector(n))&&e.push(l);for(i=[],u=0,p=e.length;u<p;u++)if(n=e[u],i.unshift(n),o=i.join(" > "),this.testSelector(t,o))return o;return null},("undefined"!=typeof exports&&null!==exports?exports:this).CssSelectorGenerator=t}).call(this);
  var cssSelectorGenerator = new CssSelectorGenerator();
  /**
   * @popperjs/core v2.3.3 - MIT License
   */
  !function(e,t){"object"==typeof exports&&"undefined"!=typeof module?t(exports):"function"==typeof define&&define.amd?define(["exports"],t):t((e=e||self).Popper={})}(this,(function(e){function t(e){return{width:(e=e.getBoundingClientRect()).width,height:e.height,top:e.top,right:e.right,bottom:e.bottom,left:e.left,x:e.left,y:e.top}}function r(e){return"[object Window]"!==e.toString()?(e=e.ownerDocument)?e.defaultView:window:e}function n(e){return{scrollLeft:(e=r(e)).pageXOffset,scrollTop:e.pageYOffset}}function o(e){return e instanceof r(e).Element||e instanceof Element}function i(e){return e instanceof r(e).HTMLElement||e instanceof HTMLElement}function a(e){return e?(e.nodeName||"").toLowerCase():null}function s(e){return(o(e)?e.ownerDocument:e.document).documentElement}function f(e){return t(s(e)).left+n(e).scrollLeft}function p(e,o,p){void 0===p&&(p=!1),e=t(e);var c={scrollLeft:0,scrollTop:0},u={x:0,y:0};return p||("body"!==a(o)&&(c=o!==r(o)&&i(o)?{scrollLeft:o.scrollLeft,scrollTop:o.scrollTop}:n(o)),i(o)?((u=t(o)).x+=o.clientLeft,u.y+=o.clientTop):(o=s(o))&&(u.x=f(o))),{x:e.left+c.scrollLeft-u.x,y:e.top+c.scrollTop-u.y,width:e.width,height:e.height}}function c(e){return{x:e.offsetLeft,y:e.offsetTop,width:e.offsetWidth,height:e.offsetHeight}}function u(e){return"html"===a(e)?e:e.assignedSlot||e.parentNode||e.host||s(e)}function l(e){return r(e).getComputedStyle(e)}function d(e,t){void 0===t&&(t=[]);var n=function e(t){if(0<=["html","body","#document"].indexOf(a(t)))return t.ownerDocument.body;if(i(t)){var r=l(t);if(/auto|scroll|overlay|hidden/.test(r.overflow+r.overflowY+r.overflowX))return t}return e(u(t))}(e);e="body"===a(n);var o=r(n);return n=e?[o].concat(o.visualViewport||[]):n,t=t.concat(n),e?t:t.concat(d(u(n)))}function m(e){return i(e)&&"fixed"!==l(e).position?e.offsetParent:null}function h(e){var t=r(e);for(e=m(e);e&&0<=["table","td","th"].indexOf(a(e));)e=m(e);return e&&"body"===a(e)&&"static"===l(e).position?t:e||t}function v(e){var t=new Map,r=new Set,n=[];return e.forEach((function(e){t.set(e.name,e)})),e.forEach((function(e){r.has(e.name)||function e(o){r.add(o.name),[].concat(o.requires||[],o.requiresIfExists||[]).forEach((function(n){r.has(n)||(n=t.get(n))&&e(n)})),n.push(o)}(e)})),n}function g(e){var t;return function(){return t||(t=new Promise((function(r){Promise.resolve().then((function(){t=void 0,r(e())}))}))),t}}function b(e){return e.split("-")[0]}function y(){for(var e=arguments.length,t=Array(e),r=0;r<e;r++)t[r]=arguments[r];return!t.some((function(e){return!(e&&"function"==typeof e.getBoundingClientRect)}))}function w(e){void 0===e&&(e={});var t=e.defaultModifiers,r=void 0===t?[]:t,n=void 0===(e=e.defaultOptions)?F:e;return function(e,t,i){function a(){f.forEach((function(e){return e()})),f=[]}void 0===i&&(i=n);var s={placement:"bottom",orderedModifiers:[],options:Object.assign({},F,{},n),modifiersData:{},elements:{reference:e,popper:t},attributes:{},styles:{}},f=[],u=!1,l={state:s,setOptions:function(i){return a(),s.options=Object.assign({},n,{},s.options,{},i),s.scrollParents={reference:o(e)?d(e):e.contextElement?d(e.contextElement):[],popper:d(t)},i=function(e){var t=v(e);return C.reduce((function(e,r){return e.concat(t.filter((function(e){return e.phase===r})))}),[])}(function(e){var t=e.reduce((function(e,t){var r=e[t.name];return e[t.name]=r?Object.assign({},r,{},t,{options:Object.assign({},r.options,{},t.options),data:Object.assign({},r.data,{},t.data)}):t,e}),{});return Object.keys(t).map((function(e){return t[e]}))}([].concat(r,s.options.modifiers))),s.orderedModifiers=i.filter((function(e){return e.enabled})),s.orderedModifiers.forEach((function(e){var t=e.name,r=e.options;r=void 0===r?{}:r,"function"==typeof(e=e.effect)&&(t=e({state:s,name:t,instance:l,options:r}),f.push(t||function(){}))})),l.update()},forceUpdate:function(){if(!u){var e=s.elements,t=e.reference;if(y(t,e=e.popper))for(s.rects={reference:p(t,h(e),"fixed"===s.options.strategy),popper:c(e)},s.reset=!1,s.placement=s.options.placement,s.orderedModifiers.forEach((function(e){return s.modifiersData[e.name]=Object.assign({},e.data)})),t=0;t<s.orderedModifiers.length;t++)if(!0===s.reset)s.reset=!1,t=-1;else{var r=s.orderedModifiers[t];e=r.fn;var n=r.options;n=void 0===n?{}:n,r=r.name,"function"==typeof e&&(s=e({state:s,options:n,name:r,instance:l})||s)}}},update:g((function(){return new Promise((function(e){l.forceUpdate(),e(s)}))})),destroy:function(){a(),u=!0}};return y(e,t)?(l.setOptions(i).then((function(e){!u&&i.onFirstUpdate&&i.onFirstUpdate(e)})),l):l}}function x(e){return 0<=["top","bottom"].indexOf(e)?"x":"y"}function O(e){var t=e.reference,r=e.element,n=(e=e.placement)?b(e):null;e=e?e.split("-")[1]:null;var o=t.x+t.width/2-r.width/2,i=t.y+t.height/2-r.height/2;switch(n){case"top":o={x:o,y:t.y-r.height};break;case"bottom":o={x:o,y:t.y+t.height};break;case"right":o={x:t.x+t.width,y:i};break;case"left":o={x:t.x-r.width,y:i};break;default:o={x:t.x,y:t.y}}if(null!=(n=n?x(n):null))switch(i="y"===n?"height":"width",e){case"start":o[n]=Math.floor(o[n])-Math.floor(t[i]/2-r[i]/2);break;case"end":o[n]=Math.floor(o[n])+Math.ceil(t[i]/2-r[i]/2)}return o}function M(e){var t,n=e.popper,o=e.popperRect,i=e.placement,a=e.offsets,f=e.position,p=e.gpuAcceleration,c=e.adaptive,u=window.devicePixelRatio||1;e=Math.round(a.x*u)/u||0,u=Math.round(a.y*u)/u||0;var l=a.hasOwnProperty("x");a=a.hasOwnProperty("y");var d,m="left",v="top",g=window;if(c){var b=h(n);b===r(n)&&(b=s(n)),"top"===i&&(v="bottom",u-=b.clientHeight-o.height,u*=p?1:-1),"left"===i&&(m="right",e-=b.clientWidth-o.width,e*=p?1:-1)}return n=Object.assign({position:f},c&&V),p?Object.assign({},n,((d={})[v]=a?"0":"",d[m]=l?"0":"",d.transform=2>(g.devicePixelRatio||1)?"translate("+e+"px, "+u+"px)":"translate3d("+e+"px, "+u+"px, 0)",d)):Object.assign({},n,((t={})[v]=a?u+"px":"",t[m]=l?e+"px":"",t.transform="",t))}function j(e){return e.replace(/left|right|bottom|top/g,(function(e){return I[e]}))}function E(e){return e.replace(/start|end/g,(function(e){return _[e]}))}function D(e,t){var r=!(!t.getRootNode||!t.getRootNode().host);if(e.contains(t))return!0;if(r)do{if(t&&e.isSameNode(t))return!0;t=t.parentNode||t.host}while(t);return!1}function P(e){return Object.assign({},e,{left:e.x,top:e.y,right:e.x+e.width,bottom:e.y+e.height})}function L(e,o){if("viewport"===o){var a=r(e);e=a.visualViewport,o=a.innerWidth,a=a.innerHeight,e&&/iPhone|iPod|iPad/.test(navigator.platform)&&(o=e.width,a=e.height),e=P({width:o,height:a,x:0,y:0})}else i(o)?e=t(o):(e=r(a=s(e)),o=n(a),(a=p(s(a),e)).height=Math.max(a.height,e.innerHeight),a.width=Math.max(a.width,e.innerWidth),a.x=-o.scrollLeft,a.y=-o.scrollTop,e=P(a));return e}function k(e,t,n){return t="clippingParents"===t?function(e){var t=d(e),r=0<=["absolute","fixed"].indexOf(l(e).position)&&i(e)?h(e):e;return o(r)?t.filter((function(e){return o(e)&&D(e,r)})):[]}(e):[].concat(t),(n=(n=[].concat(t,[n])).reduce((function(t,n){var o=L(e,n),p=r(n=i(n)?n:s(e)),c=i(n)?l(n):{};parseFloat(c.borderTopWidth);var u=parseFloat(c.borderRightWidth)||0,d=parseFloat(c.borderBottomWidth)||0,m=parseFloat(c.borderLeftWidth)||0;c="html"===a(n);var h=f(n),v=n.clientWidth+u,g=n.clientHeight+d;return c&&50<p.innerHeight-n.clientHeight&&(g=p.innerHeight-d),d=c?0:n.clientTop,u=n.clientLeft>m?u:c?p.innerWidth-v-h:n.offsetWidth-v,p=c?p.innerHeight-g:n.offsetHeight-g,n=c?h:n.clientLeft,t.top=Math.max(o.top+d,t.top),t.right=Math.min(o.right-u,t.right),t.bottom=Math.min(o.bottom-p,t.bottom),t.left=Math.max(o.left+n,t.left),t}),L(e,n[0]))).width=n.right-n.left,n.height=n.bottom-n.top,n.x=n.left,n.y=n.top,n}function B(e){return Object.assign({},{top:0,right:0,bottom:0,left:0},{},e)}function W(e,t){return t.reduce((function(t,r){return t[r]=e,t}),{})}function H(e,r){void 0===r&&(r={});var n=r;r=void 0===(r=n.placement)?e.placement:r;var i=n.boundary,a=void 0===i?"clippingParents":i,f=void 0===(i=n.rootBoundary)?"viewport":i;i=void 0===(i=n.elementContext)?"popper":i;var p=n.altBoundary,c=void 0!==p&&p;n=B("number"!=typeof(n=void 0===(n=n.padding)?0:n)?n:W(n,R));var u=e.elements.reference;p=e.rects.popper,a=k(o(c=e.elements[c?"popper"===i?"reference":"popper":i])?c:c.contextElement||s(e.elements.popper),a,f),c=O({reference:f=t(u),element:p,strategy:"absolute",placement:r}),p=P(Object.assign({},p,{},c)),f="popper"===i?p:f;var l={top:a.top-f.top+n.top,bottom:f.bottom-a.bottom+n.bottom,left:a.left-f.left+n.left,right:f.right-a.right+n.right};if(e=e.modifiersData.offset,"popper"===i&&e){var d=e[r];Object.keys(l).forEach((function(e){var t=0<=["right","bottom"].indexOf(e)?1:-1,r=0<=["top","bottom"].indexOf(e)?"y":"x";l[e]+=d[r]*t}))}return l}function T(e,t,r){return void 0===r&&(r={x:0,y:0}),{top:e.top-t.height-r.y,right:e.right-t.width+r.x,bottom:e.bottom-t.height+r.y,left:e.left-t.width-r.x}}function A(e){return["top","right","bottom","left"].some((function(t){return 0<=e[t]}))}var R=["top","bottom","right","left"],q=R.reduce((function(e,t){return e.concat([t+"-start",t+"-end"])}),[]),S=[].concat(R,["auto"]).reduce((function(e,t){return e.concat([t,t+"-start",t+"-end"])}),[]),C="beforeRead read afterRead beforeMain main afterMain beforeWrite write afterWrite".split(" "),F={placement:"bottom",modifiers:[],strategy:"absolute"},N={passive:!0},V={top:"auto",right:"auto",bottom:"auto",left:"auto"},I={left:"right",right:"left",bottom:"top",top:"bottom"},_={start:"end",end:"start"},U=[{name:"eventListeners",enabled:!0,phase:"write",fn:function(){},effect:function(e){var t=e.state,n=e.instance,o=(e=e.options).scroll,i=void 0===o||o,a=void 0===(e=e.resize)||e,s=r(t.elements.popper),f=[].concat(t.scrollParents.reference,t.scrollParents.popper);return i&&f.forEach((function(e){e.addEventListener("scroll",n.update,N)})),a&&s.addEventListener("resize",n.update,N),function(){i&&f.forEach((function(e){e.removeEventListener("scroll",n.update,N)})),a&&s.removeEventListener("resize",n.update,N)}},data:{}},{name:"popperOffsets",enabled:!0,phase:"read",fn:function(e){var t=e.state;t.modifiersData[e.name]=O({reference:t.rects.reference,element:t.rects.popper,strategy:"absolute",placement:t.placement})},data:{}},{name:"computeStyles",enabled:!0,phase:"beforeWrite",fn:function(e){var t=e.state,r=e.options;e=void 0===(e=r.gpuAcceleration)||e,r=void 0===(r=r.adaptive)||r,e={placement:b(t.placement),popper:t.elements.popper,popperRect:t.rects.popper,gpuAcceleration:e},null!=t.modifiersData.popperOffsets&&(t.styles.popper=Object.assign({},t.styles.popper,{},M(Object.assign({},e,{offsets:t.modifiersData.popperOffsets,position:t.options.strategy,adaptive:r})))),null!=t.modifiersData.arrow&&(t.styles.arrow=Object.assign({},t.styles.arrow,{},M(Object.assign({},e,{offsets:t.modifiersData.arrow,position:"absolute",adaptive:!1})))),t.attributes.popper=Object.assign({},t.attributes.popper,{"data-popper-placement":t.placement})},data:{}},{name:"applyStyles",enabled:!0,phase:"write",fn:function(e){var t=e.state;Object.keys(t.elements).forEach((function(e){var r=t.styles[e]||{},n=t.attributes[e]||{},o=t.elements[e];i(o)&&a(o)&&(Object.assign(o.style,r),Object.keys(n).forEach((function(e){var t=n[e];!1===t?o.removeAttribute(e):o.setAttribute(e,!0===t?"":t)})))}))},effect:function(e){var t=e.state,r={popper:{position:t.options.strategy,left:"0",top:"0",margin:"0"},arrow:{position:"absolute"},reference:{}};return Object.assign(t.elements.popper.style,r.popper),t.elements.arrow&&Object.assign(t.elements.arrow.style,r.arrow),function(){Object.keys(t.elements).forEach((function(e){var n=t.elements[e],o=t.attributes[e]||{};e=Object.keys(t.styles.hasOwnProperty(e)?t.styles[e]:r[e]).reduce((function(e,t){return e[t]="",e}),{}),i(n)&&a(n)&&(Object.assign(n.style,e),Object.keys(o).forEach((function(e){n.removeAttribute(e)})))}))}},requires:["computeStyles"]},{name:"offset",enabled:!0,phase:"main",requires:["popperOffsets"],fn:function(e){var t=e.state,r=e.name,n=void 0===(e=e.options.offset)?[0,0]:e,o=(e=S.reduce((function(e,r){var o=t.rects,i=b(r),a=0<=["left","top"].indexOf(i)?-1:1,s="function"==typeof n?n(Object.assign({},o,{placement:r})):n;return o=(o=s[0])||0,s=((s=s[1])||0)*a,i=0<=["left","right"].indexOf(i)?{x:s,y:o}:{x:o,y:s},e[r]=i,e}),{}))[t.placement],i=o.x;o=o.y,null!=t.modifiersData.popperOffsets&&(t.modifiersData.popperOffsets.x+=i,t.modifiersData.popperOffsets.y+=o),t.modifiersData[r]=e}},{name:"flip",enabled:!0,phase:"main",fn:function(e){var t=e.state,r=e.options;if(e=e.name,!t.modifiersData[e]._skip){var n=r.fallbackPlacements,o=r.padding,i=r.boundary,a=r.rootBoundary,s=r.altBoundary,f=r.flipVariations,p=void 0===f||f,c=r.allowedAutoPlacements;f=b(r=t.options.placement),n=n||(f!==r&&p?function(e){if("auto"===b(e))return[];var t=j(e);return[E(e),t,E(t)]}(r):[j(r)]);var u=[r].concat(n).reduce((function(e,r){return e.concat("auto"===b(r)?function(e,t){void 0===t&&(t={});var r=t.boundary,n=t.rootBoundary,o=t.padding,i=t.flipVariations,a=t.allowedAutoPlacements,s=void 0===a?S:a,f=t.placement.split("-")[1],p=(f?i?q:q.filter((function(e){return e.split("-")[1]===f})):R).filter((function(e){return 0<=s.indexOf(e)})).reduce((function(t,i){return t[i]=H(e,{placement:i,boundary:r,rootBoundary:n,padding:o})[b(i)],t}),{});return Object.keys(p).sort((function(e,t){return p[e]-p[t]}))}(t,{placement:r,boundary:i,rootBoundary:a,padding:o,flipVariations:p,allowedAutoPlacements:c}):r)}),[]);n=t.rects.reference,r=t.rects.popper;var l=new Map;f=!0;for(var d=u[0],m=0;m<u.length;m++){var h=u[m],v=b(h),g="start"===h.split("-")[1],y=0<=["top","bottom"].indexOf(v),w=y?"width":"height",x=H(t,{placement:h,boundary:i,rootBoundary:a,altBoundary:s,padding:o});if(g=y?g?"right":"left":g?"bottom":"top",n[w]>r[w]&&(g=j(g)),w=j(g),(v=[0>=x[v],0>=x[g],0>=x[w]]).every((function(e){return e}))){d=h,f=!1;break}l.set(h,v)}if(f)for(s=function(e){var t=u.find((function(t){if(t=l.get(t))return t.slice(0,e).every((function(e){return e}))}));if(t)return d=t,"break"},n=p?3:1;0<n&&"break"!==s(n);n--);t.placement!==d&&(t.modifiersData[e]._skip=!0,t.placement=d,t.reset=!0)}},requiresIfExists:["offset"],data:{_skip:!1}},{name:"preventOverflow",enabled:!0,phase:"main",fn:function(e){var t=e.state,r=e.options;e=e.name;var n=r.mainAxis,o=void 0===n||n;n=void 0!==(n=r.altAxis)&&n;var i=r.tether;i=void 0===i||i;var a=r.tetherOffset,s=void 0===a?0:a;r=H(t,{boundary:r.boundary,rootBoundary:r.rootBoundary,padding:r.padding,altBoundary:r.altBoundary}),a=b(t.placement);var f=t.placement.split("-")[1],p=!f,u=x(a);a="x"===u?"y":"x";var l=t.modifiersData.popperOffsets,d=t.rects.reference,m=t.rects.popper,v="function"==typeof s?s(Object.assign({},t.rects,{placement:t.placement})):s;if(s={x:0,y:0},l){if(o){var g="y"===u?"top":"left",y="y"===u?"bottom":"right",w="y"===u?"height":"width";o=l[u];var O=l[u]+r[g],M=l[u]-r[y],j=i?-m[w]/2:0,E="start"===f?d[w]:m[w];f="start"===f?-m[w]:-d[w],m=t.elements.arrow,m=i&&m?c(m):{width:0,height:0};var D=t.modifiersData["arrow#persistent"]?t.modifiersData["arrow#persistent"].padding:{top:0,right:0,bottom:0,left:0};g=D[g],y=D[y],m=Math.max(0,Math.min(d[w],m[w])),E=p?d[w]/2-j-m-g-v:E-m-g-v,p=p?-d[w]/2+j+m+y+v:f+m+y+v,v=t.elements.arrow&&h(t.elements.arrow),d=t.modifiersData.offset?t.modifiersData.offset[t.placement][u]:0,v=l[u]+E-d-(v?"y"===u?v.clientTop||0:v.clientLeft||0:0),p=l[u]+p-d,i=Math.max(i?Math.min(O,v):O,Math.min(o,i?Math.max(M,p):M)),l[u]=i,s[u]=i-o}n&&(n=l[a],i=Math.max(n+r["x"===u?"top":"left"],Math.min(n,n-r["x"===u?"bottom":"right"])),l[a]=i,s[a]=i-n),t.modifiersData[e]=s}},requiresIfExists:["offset"]},{name:"arrow",enabled:!0,phase:"main",fn:function(e){var t,r=e.state;e=e.name;var n=r.elements.arrow,o=r.modifiersData.popperOffsets,i=b(r.placement),a=x(i);if(i=0<=["left","right"].indexOf(i)?"height":"width",n&&o){var s=r.modifiersData[e+"#persistent"].padding;n=c(n);var f="y"===a?"top":"left",p="y"===a?"bottom":"right",u=r.rects.reference[i]+r.rects.reference[a]-o[a]-r.rects.popper[i];o=o[a]-r.rects.reference[a];var l=r.elements.arrow&&h(r.elements.arrow);u=(l=l?"y"===a?l.clientHeight||0:l.clientWidth||0:0)/2-n[i]/2+(u/2-o/2),i=Math.max(s[f],Math.min(u,l-n[i]-s[p])),r.modifiersData[e]=((t={})[a]=i,t.centerOffset=i-u,t)}},effect:function(e){var t=e.state,r=e.options;e=e.name;var n=r.element;if(n=void 0===n?"[data-popper-arrow]":n,r=void 0===(r=r.padding)?0:r,null!=n){if("string"==typeof n&&!(n=t.elements.popper.querySelector(n)))return;D(t.elements.popper,n)&&(t.elements.arrow=n,t.modifiersData[e+"#persistent"]={padding:B("number"!=typeof r?r:W(r,R))})}},requires:["popperOffsets"],requiresIfExists:["preventOverflow"]},{name:"hide",enabled:!0,phase:"main",requiresIfExists:["preventOverflow"],fn:function(e){var t=e.state;e=e.name;var r=t.rects.reference,n=t.rects.popper,o=t.modifiersData.preventOverflow,i=H(t,{elementContext:"reference"}),a=H(t,{altBoundary:!0});r=T(i,r),n=T(a,n,o),o=A(r),a=A(n),t.modifiersData[e]={referenceClippingOffsets:r,popperEscapeOffsets:n,isReferenceHidden:o,hasPopperEscaped:a},t.attributes.popper=Object.assign({},t.attributes.popper,{"data-popper-reference-hidden":o,"data-popper-escaped":a})}}],z=w({defaultModifiers:U});e.createPopper=z,e.defaultModifiers=U,e.detectOverflow=H,e.popperGenerator=w,Object.defineProperty(e,"__esModule",{value:!0})}));

  var all = document.getElementsByTagName("*");
  var beforeE = null;

  $("*").unbind();

  for (var i=0, max=all.length; i < max; i++) {
    var tarE = all[i];
    if(tarE.tagName == "body" || tarE.tagName == "html") { continue; }

    tarE.addEventListener('click', function(e) {
      var tE = e.target;
      if(tE.classList.contains("631181e4fd098a70d1f99c00192da69f9e5748935f85479e76786cc74d2772bb") || tE == beforeE) { return; };

      if(beforeE) {
        beforeE.classList.remove("coloring-depth");
      }

      beforeE = tE;

      var message = "";
			var selector = cssSelectorGenerator.getSelector(beforeE);
      message += "<p>" + cssSelectorGenerator.getSelector(beforeE) + "</p>";
      message += "<button onclick=\"window.location='/decide?selector="+encodeURI(selector)+"\" class='631181e4fd098a70d1f99c00192da69f9e5748935f85479e76786cc74d2772bb'>決定</button>";

      const tooltip = document.querySelector('#tooltip-7c392fb14fe25f428f3194f59b5b01e1c6adf8702e41755abb774812de3238dc');
      tooltip.innerHTML = message;

      beforeE.classList.add("coloring-depth");

      Popper.createPopper(beforeE, tooltip, {
        placement: 'bottom',
      });
    });
  };

  document.body.innerHTML += ` + "`" + `
    <div id="tooltip-7c392fb14fe25f428f3194f59b5b01e1c6adf8702e41755abb774812de3238dc" class="631181e4fd098a70d1f99c00192da69f9e5748935f85479e76786cc74d2772bb" role="tooltip"></div>
  ` + "`" + `;
});
</script>
	`
}
