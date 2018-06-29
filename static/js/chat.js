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
            user: 0,
            room_id: 0,
            current_room_id: 0,
            current_room: {},
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
            current_user_focus: {},
            navbar_type: 'room_list',
            room_list_unread_count: 0
        },
        created: function() {
            this.get_room_list(true,this.search_data);
            this.websocket_start();
        },
        methods: {
            // get_user: function () {
            //     var that = this;
            //     $.ajax("/api/user", {
            //         data: {},
            //         dataType: 'json',
            //         type: 'GET',
            //         success: function(response) {
            //             that.user = parseInt(response);
            //         }
            //     });
            // },
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
                        that.rooms = response===null ? [] : response;
                        if(that.rooms.length>0 && is_init) {
                            that.current_room_id = that.current_room_id>0 ? that.current_room_id : that.rooms[0].rid;
                            if(app.getQueryString("room_id")!=-1){
                                /*url room id*/
                                that.current_room_id = app.getQueryString("room_id");
                            }
                            that.changeRoom(that.current_room_id);
                            that.initRoomItemContent();
                        }
                        that.calc_roomlist_unread_count();
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
                        that.user_focus = response===null ? [] : response;
                        console.log(that.user_focus);
                        if(that.user_focus.length>0) {
                            that.current_user_focus = that.user_focus[0];
                        }
                    }
                });
            },
            changeUserFocus: function(focusUser) {
                app.current_user_focus = focusUser;
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
                        console.log(response);
                        location.href = "?room_id=" + response.rid;
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
                        if (response === null)
                        {
                            that.messages = []
                        } else {
                            that.messages = response;
                            app.scrollToEnd();
                        }
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
            submitContent: function(){
                if($("#contentBox").html().replace(/(^\s*)|(\s*$)/g, "")=="") {
                    toastr.warning("请输入数据");
                    return;
                }
                var reg = /[!！@#$%^*]/g;
                if(reg.test($("#contentBox").html())){
                    toastr.warning('输入内容不能包含特殊字符！');
                    return;
                }
                var imgLen = $("#contentBox").html().match(/<img[^>]+>/g);
                if(imgLen != null){
                    if(imgLen.length >10 ){
                        toastr.warning("上传表情不能超过10个！");
                        return false;
                    }
                }
                var content = $("#contentBox").html().replace(/<[^>]+>/g,"");//去掉所有的html标记
                if(content.length>400) {
                    toastr.warning('输入内容不能查超出400字符！');
                    return;
                }

                app.submitCommand({method: 'send_message', data: $("#contentBox").html()});
                app.scrollToEnd();
            },
            submitCommand: function(msg) {
                console.log(msg);
                if (this.conn.readyState===1) {
                    this.conn.send(objToJson(msg));
                }

            },
            deleteRoom: function(room_id) {
                if (this.conn.readyState===1) {
                    var msg = {
                        method: 'delete_room',
                        data:{
                            room_id: room_id
                        }
                    };
                    console.log(msg);
                    this.conn.send(objToJson(msg));
                }
            },
            websocket_start: function () {
                if (window["WebSocket"]) {
                    var that = this;
                    this.conn = new WebSocket("ws://" + document.location.host + "/ws");
                    this.conn.onmessage = function (evt) {
                        var messages = evt.data.split('\n');
                        console.log("receive:" + messages[0]);
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
                            default:break;
                        }
                    };

                    this.conn.onopen = function () {
                        console.log("socket has been opened");
                        var message = {
                            method: "test_connect",
                            data:"200"
                        };
                        that.conn.send(objToJson(message));
                    };
                } else {
                    var item = document.createElement("div");
                    item.innerHTML = "<b class='font-size12'>Your browser does not support WebSockets.</b>";
                    this.appendLog(item);
                }
            },
            handleTestConnect: function(msg) {
                app.has_connected = true;
                app.user = parseInt(msg.data[0]);
            },
            handleSendMessage:function(msg) {
                /*添加时间提示*/
                // if(1) {
                //     var dateinfo = {
                //                         "isDate": true,
                //                         "create_date": msg.data[0].create_date,
                //     };
                //     this.messages.push(dateinfo);
                // }
                this.messages.push(msg.data[0]);
                this.calc_roomlist_unread_count();
                this.reformRoomUnread(msg.data[0], true);
                dismissContainer();
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
                        app.changeRoom(app.current_room_id);
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
            changeRoom: function (room_id) {
                this.messages = [];
                if(app.rooms.length>0) {
                    this.get_room_messages(room_id);
                    this.submitCommand({
                        method: 'enter_room',
                        data:{
                            room_id: room_id
                        }
                    });
                    app.reformActiveRoom(room_id); //init active room
                }

            },
            searchSubStr: function(str,subStr){
                var pos = str.indexOf(subStr);
                while(pos>-1){
                    app.positions.push(pos);
                    pos = str.indexOf(subStr,pos+1);
                }
            },
            scrollToEnd:function(){//滚动到底部
                $('#log').animate({scrollTop:$('#log').height() + 99999 + 'px'},100);
                $("#contentBox").html("");
            },
            reformActiveRoom: function(room_id) {
                if(app.rooms.length>0) {
                    for(var i=0; i< app.rooms.length; i++) {
                        if(parseInt(app.rooms[i].rid)==room_id) {
                            Vue.set(app.rooms[i], 'isActive', true);
                            app.current_room_id = room_id;
                            Vue.set(app.rooms[i], 'unread', 0);
                            app.current_room = app.rooms[i];
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
            }
        }
    });