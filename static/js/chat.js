function jsonToObj(str) {
        var obj = JSON.parse(str);
        return obj;
}
function objToJson(obj) {
    var str = JSON.stringify(obj);
    return str;
}
function chatChooseEmoji() {
    $(".emoji-container").css({'display' : 'inline-block'});
}
function dismissContainer() {
    $(".emoji-container").css({'display' : 'none'});
}
function chatInertEmoji(e) {
    var event = window.event || e;
    var event_target = event.srcElement ? event.srcElement : event.target;
    var emoji_src = $(event_target).attr('src');
    $("#contentBox").append("<img src='" + emoji_src + "' model_id='0' class='topic-emojis'>");
}
var emojis = new Array();
var max_length = 1;
for(var i=0; i<5; i++) {
    emojis[i] = new Array();
    for(var j=0; j<15; j++) {
        emojis[i][j] = max_length;
        max_length++;
    }
}

var app = new Vue({
        el: '#app',
        data: {
            isMobile: 'yes',
            user: 0,
            selfInfo: {},
            room_id: 0,
            current_room_id: 0,
            current_room: {
                index: -1,
                images: '',
                des: ''
            },
            rooms: [],
            conn: false,
            has_connected: false,
            messages: [],
            positions:[],
            username: "",
            emojiss: emojis,
            image_prefix: "http://cdn.ggac.net/",
            search_data: "",
            user_focus: [],
            current_user_focus: {
                id: 0,
                avatar_image_small: '',
                username: '',
                num_focus: 0,
                num_focused: 0,
                is_focused: true
            },
            navbar_type: 'room_list',
            room_list_unread_count: 0,
            desc: "", //已输入字数
            remnant: 400 //限定字数
        },
        created: function() {
            this.websocket_start();
            this.initData();
            this.get_user();
            this.get_room_list(true,this.search_data);
        },
        methods: {
            initData:function(){
                var width = document.body.clientWidth;
                if(width>=1200) {
                    this.isMobile = 'no';
                } else {
                    this.isMobile = 'yes';
                }
            },
            descInput: function() {
                var txtVal = this.desc.length;
                this.remnant = 400 - txtVal;
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
                        } else {
                            toastr.warning(response.result);
                        }
                    }
                });
            },
            get_room_list: function (is_init,value) {
                console.log(this.navbar_type);
                var that = this;
                this.search_data = is_init ? "" : this.search_data;
                that.navbar_type = "room_list";
                $.ajax("/api/room/list", {
                    data: {
                        key: value
                    },
                    dataType: 'json',
                    type: 'GET',
                    success: function(response) {
                        if(response.return_code==0) {
                            that.rooms = response.data;
                        } else {
                            toastr.warning(response.result);
                        }
                        if(that.rooms.length>0) {
                            if(is_init) {
                                that.current_room_id = that.current_room_id>0 ? that.current_room_id : that.rooms[0].rid;
                                if(app.getQueryString("room_id")!=-1){
                                    /*url room id*/
                                    that.current_room_id = app.getQueryString("room_id");
                                }
                            } else {
                                //通讯录切换回私信，默认第一个
                                that.current_room_id = that.rooms[0].rid;
                            }

                            that.changeRoom(that.current_room_id, true);
                            that.initRoomItemContent();
                        }
                        that.calc_roomlist_unread_count();

                        Vue.nextTick(function(){
                            $('#roomList').mCustomScrollbar("destroy");
                            $('#roomList').mCustomScrollbar("");
                        });
                    }
                });
            },
            get_user_focus_list: function(value) {

                var that = this;
                this.search_data = "";
                that.navbar_type = "focus_list";
                console.log(this.navbar_type);
                $.ajax("/api/user/focused/list", {
                    data: {
                        key: value
                    },
                    dataType: 'json',
                    type: 'GET',
                    success: function(response) {
                        if(response.return_code==0) {
                            that.user_focus = response.data;
                        } else {
                            toastr.warning(response.result);
                        }
                        if(that.user_focus.length>0) {
                            that.current_user_focus = that.user_focus[0];
                            that.getFocusUserInfo(that.current_user_focus.id); //拉取当前通讯联系人详情数据（粉丝、关注）
                        }
                        Vue.nextTick(function(){
                            $('#focusList').mCustomScrollbar("destroy");
                            $('#focusList').mCustomScrollbar();
                        });
                    }
                });
            },
            getFocusUserInfo: function(user_id) {
                var that = this;
                $.ajax("/api/user/detail", {
                    data: {
                        user_id: user_id
                    },
                    dataType: 'json',
                    type: 'GET',
                    success: function(response) {
                        if(response.return_code==0) {
                            Vue.set(app.current_user_focus, 'num_focus', response.data[0].num_focus);
                            Vue.set(app.current_user_focus, 'num_focused', response.data[0].num_focused);
                            Vue.set(app.current_user_focus, 'is_focused', true);
                        } else {
                            Vue.set(app.current_user_focus, 'num_focus', 0);
                            Vue.set(app.current_user_focus, 'num_focused', 0);
                            Vue.set(app.current_user_focus, 'is_focused', false);
                            toastr.warning(response.result);
                        }
                    }
                });
            },
            cancelFocusUser: function(user_id) {
                var that = this;
                $.ajax("/api/user/focused/cancel ", {
                    data: {
                        user_id: user_id
                    },
                    dataType: 'json',
                    type: 'POST',
                    success: function(response) {
                        toastr.success("取消关注成功");
                        Vue.set(app.current_user_focus, 'is_focused', false);
                        //app.get_user_focus_list(""); //重新拉取通讯录列表
                    }
                });
            },
            focusUser: function(user_id) {
                var that = this;
                $.ajax("/api/user/focused/agree", {
                    data: {
                        user_id: user_id
                    },
                    dataType: 'json',
                    type: 'POST',
                    success: function(response) {
                        if(response.return_code==0) {
                            toastr.success("关注成功");
                            Vue.set(app.current_user_focus, 'is_focused', true);
                        } else {
                            toastr.success(response.result);
                        }
                        //app.get_user_focus_list(""); //重新拉取通讯录列表
                    }
                });
            },
            changeUserFocus: function(focusUser) {
                app.current_user_focus = focusUser;
                app.getFocusUserInfo(app.current_user_focus.id);  //拉取当前通讯联系人详情数据（粉丝、关注）
            },
            calc_roomlist_unread_count: function() {
                app.room_list_unread_count = 0;
                for(var i=0; i<app.rooms.length; i++) {
                    var unread_num = parseInt(app.rooms[i].unread);
                    if(unread_num>0) {
                        app.room_list_unread_count += unread_num;
                    }
                }
            },
            sendUserFocus: function() {
                var that = this;
                console.log("focusUser id: " + that.current_user_focus.id);
                $.ajax("/api/room/create", {
                    data: {
                        user_ids: [that.current_user_focus.id]
                    },
                    dataType: 'json',
                    type: 'POST',
                    success: function(response) {
                        if(response.return_code==0) {
                            location.href = "?room_id=" + response.data[0].rid;
                        } else {
                            toastr.warning(response.result);
                        }
                    }
                });
            },
            get_room_messages: function (room_id) {
                var that = this;
                $.ajax("/api/room/message/list", {
                    data: {
                        'room_id': room_id
                    },
                    dataType: 'json',
                    type: 'GET',
                    success: function(response) {
                        if(response.return_code==0){
                            that.messages = response.data;
                            app.resetContentBox(); //reset send message box
                            app.scrollToEnd();
                        } else {
                            that.messages = [];
                            toastr.warning(response.result);
                        }

                        Vue.nextTick(function(){
                            $('#log').mCustomScrollbar("destroy");
                            $('#log').mCustomScrollbar("");
                            $('#log').mCustomScrollbar("scrollTo", "bottom", {scrollInertia: 0});
                        });

                       //初始化滚动条样式
                    }
                });
            },
            submitSearch: function() {
                if(this.navbar_type=='room_list') {
                    this.get_room_list(false,this.search_data);
                } else {
                    this.get_user_focus_list(this.search_data);
                }

            },
            submitContent: function(evt){

                var removeBr = $("#contentBox").val().replace(/<br>/ig,"");   //remove <br>
                var removerBrSpace = removeBr.replace(/\s+/g,"");  //remove blank space
                var removerBrSpaceNbsp = removerBrSpace.replace(/&nbsp;/ig, "");  //remove &nbsp;

                if(removeBr=="" || removerBrSpace=="" || removerBrSpaceNbsp=="") {
                    toastr.warning("请输入数据");
                    app.resetContentBox(); //reset send message box
                    evt.preventDefault();
                    return false;
                }

                var reg = /[!！@#$%^*]/g;
                if(reg.test($("#contentBox").val())){
                    toastr.warning('输入内容不能包含特殊字符！');
                    return;
                }
                var imgLen = $("#contentBox").val().match(/<img[^>]+>/g);
                if(imgLen != null){
                    if(imgLen.length >10 ){
                        toastr.warning("上传表情不能超过10个！");
                        return false;
                    }
                }
                var content = $("#contentBox").val().replace(/<[^>]+>/g,"");//去掉所有的html标记
                if(content.length>400) {
                    toastr.warning('输入内容不能超出400字符！');
                    return;
                }

                app.submitCommand({method: 'send_message', data: $("#contentBox").val()});
                app.scrollToEnd();
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
            deleteRoom: function(room_id) {
                if (this.conn.readyState===1) {
                    var msg = {
                        method: 'delete_room',
                        data:{
                            room_id: room_id
                        }
                    };
                    this.conn.send(objToJson(msg));
                }
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
                        switch (msg.method) {
                            case 'test_connect':
                                that.handleTestConnect(msg);
                                break;
                            case 'message_send':
                                that.handleSendMessage(msg);
                                break;
                            case 'room_created':
                                that.handleCreateRoom(msg);
                                break;
                            case 'unread_room':
                                that.handleUnreadRoom(msg);
                                break;
                            case 'room_deleted':
                                that.handleRoomDelete(msg);
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
            handleSendMessage:function(msg) {
                if(parseInt(app.user)==parseInt(msg.data[0].from_uid)) {
                    var  msgBox = '<div>'+
                        '<div class="self-box">' +
                        '<div class="content-date">'+
                        '<span class="font-size12 content-time">' + app.reformCreateDate(msg.data[0].create_date.split('.')[0]) + '</span>'+
                        '<span class="font-size14 content-body">' + msg.data[0].content + '</span>' +
                        '</div>' +
                        '<span class="content-header">'+
                        '<img src="' + app.image_prefix + msg.data[0].image + '" class="user-avatar mCS_img_loaded">' +
                        '</span>'+
                        '</div>' +
                        '</div>';
                } else {
                    var msgBox = '<div>' +
                                    '<div class="other-box">' +
                                        '<span class="content-header">' +
                                            '<img src="' + app.image_prefix + msg.data[0].image + '" class="user-avatar mCS_img_loaded">' +
                                        '</span>' +
                                        '<div class="content-date">' +
                                            '<span class="font-size12 content-time">' + app.reformCreateDate(msg.data[0].create_date.split('.')[0]) + '</span>' +
                                            '<span class="font-size14 content-body">' + msg.data[0].content + '</span>' +
                                        '</div>' +
                                    '</div>' +
                                '</div>';
                }


                // if(app.navbar_type == "focus_list" && app.rooms.length>0) {
                //     //当前是通讯录&接收到房间聊天记录
                //
                // } else {
                    $("#log .mCSB_container").append(msgBox);
                //}
                app.calc_roomlist_unread_count();
                app.reformRoomUnread(msg.data[0], true);

                //dismissContainer();

                if(parseInt(app.user)==parseInt(msg.data[0].from_uid)) {
                    //自己发送成功，清空输入框
                    app.resetContentBox(); //reset send message box
                    Vue.nextTick(function() {
                        if(IEVersion()==-1) {
                            app.rooms.sort(roomTimeCompare("create_date")); //reform room squence
                        } else {
                            app.rooms = IeRoomUpTop(JSON.parse(JSON.stringify(app.rooms)), app.current_room.index);
                            app.current_room.index = 0;
                        }
                        $('#roomList').mCustomScrollbar("destroy");
                        $('#roomList').mCustomScrollbar();
                    });
                }
                app.scrollToEnd(); //scroll to bottom
            },
            handleCreateRoom:function(msg) {
                //不在roomsm内，push(msg)
                console.log(msg);
                if(!app.hasRoom(msg.data[0].rid)) {
                    msg.data[0].unread = 0;
                    app.rooms.push(msg.data[0]);
                }
            },
            handleUnreadRoom:function(msg) {
                this.reformRoomUnread(msg.data[0], false);
                this.calc_roomlist_unread_count();
            },
            handleRoomDelete:function(msg) {
                console.log(msg.data[0]);
                //2、不是当前room直接删除
                app.deleteCurrentRoom(msg.data[0].room_id);
                app.calc_roomlist_unread_count();
                if(app.current_room_id==parseInt(msg.data[0].room_id)) {
                    //1、是当前room删除, 还有房间则切换第一个；没有房间不切换
                    //app.deleteRoom(msg.data[0].room_id);
                    if(app.rooms.length>0) {
                        //还有房间则切换第一个
                        app.current_room_id = app.rooms[0].rid;
                        //拉取当前room数据msg
                        app.changeRoom(app.current_room_id, true);
                        app.reformActiveRoom(app.rooms[0].rid);
                    } else {
                        //没有房间不切换
                        app.current_room_id = 0;
                        app.messages = [];
                    }
                }
            },
            appendLog: function(item) {
                var log = document.getElementById("log");
                var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
                log.appendChild(item);
                if (doScroll) {
                    log.scrollTop = log.scrollHeight - log.clientHeight;
                }
            },
            changeRoom: function (room_id,isInit) {
                this.messages = [];
                if(app.rooms.length>0) {
                    this.get_room_messages(room_id);
                    if(!isInit) {
                        $("#log .mCSB_container div:first-child").remove();
                    }
                    this.submitCommand({
                        method: 'enter_room',
                        data:{
                            room_id: room_id
                        }
                    });
                    app.reformActiveRoom(room_id); //init active room
                    app.resetContentBox(); //reset send message box

                }

            },
            resetContentBox: function() {
                //reset send message box
                app.desc = "";
                app.descInput();
                $("#contentBox").focus();
            },
            searchSubStr: function(str,subStr){
                var pos = str.indexOf(subStr);
                while(pos>-1){
                    app.positions.push(pos);
                    pos = str.indexOf(subStr,pos+1);
                }
            },
            scrollToEnd:function(){//滚动到底部
                    //Vue.nextTick(function () {
                        $('#log').mCustomScrollbar("update");
                        $('#log').mCustomScrollbar("scrollTo", "bottom", {scrollInertia: 0});
                    //});
            },
            reformActiveRoom: function(room_id) {
                if(app.rooms.length>0) {
                    for(var i=0; i< app.rooms.length; i++) {
                        if(parseInt(app.rooms[i].rid)==room_id) {
                            Vue.set(app.rooms[i], 'isActive', true);
                            app.current_room_id = room_id;
                            Vue.set(app.rooms[i], 'unread', 0);
                            app.current_room = app.rooms[i];
                            Vue.set(app.current_room, 'index', i);
                        } else {
                            Vue.set(app.rooms[i], 'isActive', false);
                        }
                    }
                    app.calc_roomlist_unread_count();
                }
            },
            reformRoomUnread: function(msg, isActive) {

                for(var i=0; i< app.rooms.length; i++) {
                    if(parseInt(app.rooms[i].rid)==msg.rid) {
                        if(!isActive) {
                            //不是当前房间更新未读数
                            var oldUnreadNumber = app.rooms[i].unread;
                            Vue.set(app.rooms[i], 'unread', parseInt(oldUnreadNumber) + 1);
                        }
                        Vue.set(app.rooms[i], 'content', msg.content.replace(/<img(.*?)>/g, "[图片]"));
                        Vue.set(app.rooms[i], 'create_date', msg.create_date);
                        return;
                    }
                }
            },
            hasRoom: function(room_id) {
                var flag = false;
                for(var i=0; i< app.rooms.length; i++) {
                    if(parseInt(app.rooms[i].rid)==room_id) {
                        flag = true;
                        break;
                    }
                }
                return flag;
            },
            deleteCurrentRoom: function(room_id) {
                for(var i=0; i< app.rooms.length; i++) {
                    if(parseInt(app.rooms[i].rid)==room_id) {
                        app.rooms.splice(i, 1);
                    }
                }
            },
            initRoomItemContent: function() {
                for(var i=0; i< app.rooms.length; i++) {
                    Vue.set(app.rooms[i], 'content', app.rooms[i].content.replace(/<img(.*?)>/g, "[图片]"));
                }
            },
            getQueryString: function(name) {
                /*get url param*/
                var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
                var r = window.location.search.substr(1).match(reg);
                if (r != null) return unescape(r[2]); return -1;
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