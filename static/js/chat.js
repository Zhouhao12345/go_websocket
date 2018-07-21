function jsonToObj(str) {
        var obj = JSON.parse(str);
        return obj;
}
function objToJson(obj) {
    var str = JSON.stringify(obj);
    return str;
}

var app = new Vue({
        el: '#app',
        data: {
            username: "",
            password: "",
            isMobile: 'yes',
            user: 0,
            selfInfo: {
                avatar_image_small: "",
                username: "",
                id: 0,
            },
            conn: false,
            has_connected: false,
            users: {},
            lockReconnect: false,
        },
        created: function() {
            this.initData();
            this.get_user();
        },
        methods: {
            showLoginFrame: function(){
                $(".loginFrame").show();
            },
            showRegisterFrame: function(){
                $(".RegisterFrame").show();
            },
            closeLoginFrame: function(){
                $(".loginFrame").hide();
            },
            closeRegisterFrame: function(){
                $(".RegisterFrame").hide();
            },
            logout: function(){
                var that = this;
                $.ajax("/api/user/logout", {
                    data: {},
                    dataType: 'json',
                    type: 'GET',
                    success: function(response) {
                        if(response.return_code==0) {
                            toastr.success("Logout Success");
                            window.location.reload();
                        } else {
                            toastr.warning(response.result);
                        }
                    }
                });
            },
            Register: function(){
                if(this.username.length===0){
                    toastr.warning("请填写用户名");
                }
                if(this.password.length===0){
                    toastr.warning("请填写密码");
                }
                var that = this;
                $.ajax("/api/user/register", {
                    data: {
                        username: that.username,
                        password: that.password,
                    },
                    dataType: 'json',
                    type: 'POST',
                    success: function(response) {
                        if(response.return_code==0) {
                            toastr.success("注册成功");
                            window.location.reload();
                        } else {
                            toastr.warning(response.result);
                        }
                    }
                });
            },
            login: function(){
                if(this.username.length===0){
                    toastr.warning("请填写用户名");
                }
                if(this.password.length===0){
                    toastr.warning("请填写密码");
                }
                var that = this;
                $.ajax("/api/user/login", {
                    data: {
                        username: that.username,
                        password: that.password,
                    },
                    dataType: 'json',
                    type: 'POST',
                    success: function(response) {
                        if(response.return_code==0) {
                            toastr.success("登录成功");
                            window.location.reload();
                        } else {
                            toastr.warning(response.result);
                        }
                    }
                });
            },
            initData:function(){
                var width = document.body.clientWidth;
                if(width>=1200) {
                    this.isMobile = 'no';
                } else {
                    this.isMobile = 'yes';
                }
            },
            get_user: function () {
                var that = this;
                $.ajax("/api/user", {
                    data: {},
                    dataType: 'json',
                    type: 'GET',
                    success: function(response) {
                        if(response.return_code==0) {
                            that.selfInfo = response.data[0];
                            that.websocket_start();
                            that.enterMap();
                        } else {
                            toastr.warning(response.result);
                        }
                    }
                });
            },
            submitCommand: function(msg) {
                var that = this;
                var interval = setInterval(function () {
                    console.log(that.conn.readyState);
                    console.log(msg);
                    if (that.conn.readyState===1) {
                        that.conn.send(objToJson(msg));
                        clearInterval(interval);
                    }
                }, 100);
                setTimeout(function() {
                    if(interval && that.conn.readyState!==1) {
                        //超时处理
                        clearInterval(interval);
                        that.handleErrorReceived("超时连接,请刷新页面重试");
                    }
                }, 10000);
            },
            reconnect: function() {
                if(app.lockReconnect) return;
                app.lockReconnect = true;
                setTimeout(function () {
                    app.get_user();
                    app.lockReconnect = false;
                }, 2000);
            },
            websocket_start: function () {
                if (window["WebSocket"]) {
                    var that = this;
                    this.conn = new WebSocket("ws://" + document.location.host + "/ws");
                    this.conn.onopen = function () {
                        console.log("socket has been opened");
                        var message = {
                            method: "test_connect",
                            data:"200"
                        };
                        that.conn.send(objToJson(message));
                    };
                    this.conn.onclose = function () {
                        that.reconnect();
                    };
                    this.conn.onmessage = function (evt) {
                        var messages = evt.data.split('\n');
                        var msg = jsonToObj(messages[0]);
                        switch (msg.method) {
                            case 'test_connect':
                                that.handleTestConnect(msg);
                                break;
                            case 'mapEnter':
                                that.handleMapEnter(msg);
                                break;
                            case 'mapInit':
                                that.handleMapInit(msg);
                                break;
                            case 'move':
                                that.handleMove(msg);
                                break;
                            case 'error_received':
                                that.handleErrorReceived(msg.data[0].error + ",请刷新页面重试");
                                break;
                            case 'mapLeave':
                                that.handleMapLeave(msg);
                                break;
                            default:break;
                        }
                    };
                    this.conn.onerror = function () {
                        that.handleErrorReceived("连接中断, 请刷新页面重试");
                    };
                } else {
                    that.handleErrorReceived("浏览器不支持WebSockets，请更换浏览器");
                }
            },
            handleErrorReceived: function(tip) {
                var that = this;
                that.has_connected = false;
                $("#userDisconnected").html("");
                Vue.nextTick(function(){
                    $("#userDisconnected").append("<span style='display: block; text-align: center;padding-top: 70px;'>" + tip +"</span>");
                });
            },
            handleTestConnect: function(msg) {
                Vue.nextTick(function(){
                    app.has_connected = true;
                    app.user = parseInt(msg.data[0]);
                });
            },
            handleMapInit: function(msg) {
                Vue.nextTick(function() {
                    var members = msg.data[0].members;
                    for(var i=0;i<members.length;i++)
                    {
                        app.users[members[i].id] = {
                            username : members[i].name,
                            image : members[i].image,
                        };
                        $("#rightBox").append(
                            "<span id='user"+members[i].id+"' style='" +
                            "font-size: 50px; " +
                            "text-align: center; " +
                            "background-color: red; " +
                            "position: relative; " +
                            "left: "+members[i].positionX+"%; " +
                            "top: "+members[i].positionY+"%;'>" + members[i].name +"</span>"
                        );
                    }
                    console.log(app.selfInfo);
                    document.onkeydown = function(){
                        var left = parseInt($("#user"+app.selfInfo.id.toString()).css("left").split("%")[0]);
                        var top = parseInt($("#user"+app.selfInfo.id.toString()).css("top").split("%")[0]);
                        var key = window.event.keyCode;
                        // left
                        if(key == 37){
                            app.move(left-1,top);
                        }
                        // top
                        if(key == 38){
                            app.move(left, top-1);
                        }
                        // right
                        if(key == 39){
                            app.move(left+1,top);
                        }
                        // down
                        if(key == 40){
                            app.move(left, top+1);
                        }
                    }
                })
            },
            handleMapEnter: function(msg) {
                Vue.nextTick(function(){
                    if(app.users[msg.data[0].user]==undefined)
                    {
                        app.users[msg.data[0].user] = {
                            username : msg.data[0].username,
                            image : msg.data[0].image,
                        };
                        $("#rightBox").append(
                            "<span id='user"+msg.data[0].user+"' style='" +
                            "font-size: 50px; " +
                            "text-align: center; " +
                            "background-color: red; " +
                            "position: relative; " +
                            "left: 50%; " +
                            "top: 50%;'>" + msg.data[0].username +"</span>"
                        );
                    } else {
                        $("#user"+msg.data[0].user).css("left", "50%");
                        $("#user"+msg.data[0].user).css("top", "50%");
                    }
                });
            },
            handleMove: function(msg) {
                Vue.nextTick(function(){
                    $("#user"+msg.data[0].user).css("left", msg.data[0].x+"%");
                    $("#user"+msg.data[0].user).css("top", msg.data[0].y+"%");
                });
            },
            handleMapLeave: function(msg) {
                Vue.nextTick(function(){
                    delete app.users[msg.data[0].user];
                    $("#user"+msg.data[0].user).remove();
                });
            },
            enterMap: function () {
                this.submitCommand({
                    method: 'enter_map',
                    data:{}
                });
            },
            move: function (left,top) {
                this.submitCommand({
                    method: 'move',
                    data:{
                        x: left.toString(),
                        y: top.toString(),
                    }
                });
            },
            handlePreview: function(file) {
                console.log("handlePreview");
            },
            handleRemove: function(file, fileList) {
                console.log("handleRemove");
            },
            handleProgress: function(event, file, fileList) {
                console.log("handleProgress");
            },
            handleSuccess: function(response, file, fileList) {
                console.log("handleSuccess");
                console.log(response);
            }
        }
    });