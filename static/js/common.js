function _toastrInitSet() {
    toastr.options = {
        closeButton: false,
        debug: false,
        progressBar: true,
        positionClass: "toast-top-center",
        onclick: null,
        showDuration: "300",
        hideDuration: "1000",
        timeOut: "2000",
        extendedTimeOut: "1000",
        showEasing: "swing",
        hideEasing: "linear",
        showMethod: "fadeIn",
        hideMethod: "fadeOut"
    };
}
_toastrInitSet();


Vue.use(VueLazyload, {
    preLoad: 1.3,
    error: "http://cdn.ggac.net/images/lazy_error.png",
    loading: "http://cdn.ggac.net/images/lazy_loading.png",
    attempt: 1
});

/*YY/MM/DD hh:mm*/
Vue.prototype.reformCreateDate = function(dateStr) {
    var dateStr = dateStr.replace(/\-/ig,"/");
    var date = new Date(dateStr);
    var year = date.getFullYear();
    var month = date.getMonth() + 1;
    var day = date.getDate();
    var hour = date.getHours();
    var minute = date.getMinutes();
    if(parseInt(month)<10) {
        month = "0" + month;
    }
    if(parseInt(day)<10) {
        day = "0" + day;
    }
    if(parseInt(hour)<10) {
        hour = "0" + hour;
    }
    if(parseInt(minute)<10) {
        minute = "0" + minute;
    }
    return (parseInt(year) - 2000) + "/" + month + "/" +  day + " " + hour + ":" + minute;
}

function roomTimeCompare(property){
    return function(a,b){
        var value1 = new Date(a[property]).getTime();
        var value2 =  new Date(b[property]).getTime();
        return  value2 -  value1;
    }
}

function IeRoomUpTop(newArr , index) {
    if(index!=0) {
        var obj = newArr[index];
        newArr.splice(index, 1);
        newArr.unshift(obj);
    }
    return newArr;
}

function IEVersion() {
    var userAgent = navigator.userAgent; //取得浏览器的userAgent字符串
    var isIE = userAgent.indexOf("compatible") > -1 && userAgent.indexOf("MSIE") > -1; //判断是否IE<11浏览器
    var isEdge = userAgent.indexOf("Edge") > -1 && !isIE; //判断是否IE的Edge浏览器
    var isIE11 = userAgent.indexOf('Trident') > -1 && userAgent.indexOf("rv:11.0") > -1;
    if(isIE) {
        var reIE = new RegExp("MSIE (\\d+\\.\\d+);");
        reIE.test(userAgent);
        var fIEVersion = parseFloat(RegExp["$1"]);
        if(fIEVersion == 7) {
            return 7;
        } else if(fIEVersion == 8) {
            return 8;
        } else if(fIEVersion == 9) {
            return 9;
        } else if(fIEVersion == 10) {
            return 10;
        } else {
            return 6;//IE版本<=7
        }
    } else if(isEdge) {
        return 'edge';//edge
    } else if(isIE11) {
        return 11; //IE11
    }else{
        return -1;//不是ie浏览器
    }
}