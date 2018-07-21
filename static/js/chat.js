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
                        that.handleErrorReceived("您与服务器失去联系, 请刷新页面重试");
                    };
                    this.conn.onmessage = function (evt) {
                        var messages = evt.data.split('\n');
                        var msg = jsonToObj(messages[0]);
                        console.log(msg);
                        switch (msg.method) {
                            case 'test_connect':
                                that.handleTestConnect(msg);
                                break;
                            case 'mapEnter':
                                that.handleMapEnter(msg);
                                break;
                            case 'move':
                                that.handleMove(msg);
                                break;
                            case 'error_received':
                                that.handleErrorReceived(msg.data[0].error + ",请刷新页面重试");
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
            handleMapEnter: function(msg) {
                Vue.nextTick(function(){
                    if(app.users[msg.data[0].user]==undefined)
                    {
                        app.users[msg.data[0].user] = {
                            username : msg.data[0].username,
                            image : msg.data[0].image,
                            positions: {
                                x: 100,
                                y: 100,
                            }
                        };
                        $("#rightBox").append(
                            "<span style='font-size: 50px; text-align: center; background-color: red'>" + msg.data[0].username +"</span>"
                        );
                    }
                });
            },
            handleMove: function(msg) {
                Vue.nextTick(function(){
                    app.users[msg["user"]]["positions"]["x"] = msg["x"];
                    app.users[msg["user"]]["positions"]["y"] = msg["y"];
                });
            },
            enterMap: function () {
                this.submitCommand({
                    method: 'enter_map',
                    data:{}
                });
            },
            move: function () {
                this.submitCommand({
                    method: 'move',
                    data:{
                        x: 100,
                        y: 100,
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