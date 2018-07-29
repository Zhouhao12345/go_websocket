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
            Phaser: Phaser,
            map: false,
            tileset: false,
            layer: false,
            player: false,
            facing: 'left',
            jumpTimer: 0,
            cursors: false,
            jumpButton: false,
            bg: false,
            player_name: false,
            game: false,
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
        mounted: function() {
            this.get_user();
        },
        methods: {
            preload: function(){
                app.game.load.tilemap('level1', '/static/image/games/starstruck/level1.json', null, app.Phaser.Tilemap.TILED_JSON);
                app.game.load.image('tiles-1', '/static/image/games/starstruck/tiles-1.png');
                app.game.load.spritesheet('dude', '/static/image/games/starstruck/dude.png', 32, 48);
                app.game.load.spritesheet('droid', '/static/image/games/starstruck/droid.png', 32, 32);
                app.game.load.image('starSmall', '/static/image/games/starstruck/star.png');
                app.game.load.image('starBig', '/static/image/games/starstruck/star2.png');
                app.game.load.image('background', '/static/image/games/starstruck/background2.png');
                app.websocket_start();
                app.enterMap();
            },

            created: function() {

                app.game.physics.startSystem(app.Phaser.Physics.ARCADE);

                app.game.stage.backgroundColor = '#000000';

                app.bg = app.game.add.tileSprite(0, 0, 800, 600, 'background');
                app.bg.fixedToCamera = true;

                app.map = app.game.add.tilemap('level1');

                app.map.addTilesetImage('tiles-1');

                app.map.setCollisionByExclusion([ 13, 14, 15, 16, 46, 47, 48, 49, 50, 51 ]);

                app.layer = app.map.createLayer('Tile Layer 1');

                //  Un-comment this on to see the collision tiles
                // layer.debug = true;

                app.layer.resizeWorld();

                app.game.physics.arcade.gravity.y = 250;

                app.cursors = app.game.input.keyboard.createCursorKeys();
                app.jumpButton = app.game.input.keyboard.addKey(app.Phaser.Keyboard.SPACEBAR);
            },

            update: function() {
                for(var id in app.users){
                    app.users[id].player_name.x = Math.floor(app.users[id].player.x - app.users[id].player.width / 2);
                    app.users[id].player_name.y = Math.floor(app.users[id].player.y - app.users[id].player.height / 2);
                    app.game.physics.arcade.collide(app.users[id].player, app.layer);
                    app.users[id].player.body.velocity.x = 0;
                }

                var self = app.users[app.selfInfo.id];
                if(self != undefined)
                {
                    if (app.cursors.left.isDown)
                    {
                        app.move(-150, 0, 0);
                        self.player.body.velocity.x = -150;
                        if (self.facing != 'left')
                        {
                            self.player.animations.play('left');
                            self.facing = 'left';
                        }
                    }
                    else if (app.cursors.right.isDown)
                    {
                        app.move(150, 0, 0);
                        self.player.body.velocity.x = 150;
                        if (self.facing != 'right')
                        {
                            self.player.animations.play('right');
                            self.facing = 'right';
                        }
                    }
                    else
                    {
                        if (self.facing != 'idle')
                        {
                            self.player.animations.stop();

                            if (self.facing == 'left')
                            {
                                self.player.frame = 0;
                            }
                            else
                            {
                                self.player.frame = 5;
                            }

                            self.facing = 'idle';
                        }
                    }
                    if (app.jumpButton.isDown && self.player.body.onFloor() && app.game.time.now > self.jumpTimer)
                    {
                        app.move(0, -250, 750);
                        self.player.body.velocity.y = -250;
                        self.jumpTimer = app.game.time.now + 750;
                    }
                }
            },
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
                Vue.nextTick(function () {
                    var Main = {
                        preload: app.preload,
                        create: app.created,
                        update: app.update,
                    };
                    app.game = new app.Phaser.Game(
                        800,
                        600,
                        Phaser.AUTO,
                        'phaser-example');
                    app.game.state.add("Main", Main);
                    app.game.state.start("Main");
                });
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
                            that.initData();
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
                        // that.reconnect();
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
                        app.userCreated(members[i].id, members[i].name);
                    }
                })
            },
            userCreated: function(id, username){
                var player = app.game.add.sprite(32, 32, 'dude');
                app.game.physics.enable(player, app.Phaser.Physics.ARCADE);

                player.body.bounce.y = 0.2;
                player.body.collideWorldBounds = true;
                player.body.setSize(20, 32, 5, 16);

                player.animations.add('left', [0, 1, 2, 3], 10, true);
                player.animations.add('turn', [4], 20, true);
                player.animations.add('right', [5, 6, 7, 8], 10, true);

                player_name = app.game.add.text(16, 16, username, { fontSize: '5px', fill: '#000' });
                if(id == app.selfInfo.id){
                    app.game.camera.follow(player);
                }
                app.users[id] = {
                    player: player,
                    player_name: player_name,
                    facing: "left",
                    jumpTimer: 0,
                }
            },
            handleMapEnter: function(msg) {
                Vue.nextTick(function(){
                    if(app.users[msg.data[0].user]==undefined)
                    {
                        app.userCreated(msg.data[0].user, msg.data[0].username);
                    } else {
                        app.users[msg.data[0].user].player.x = 32;
                        app.users[msg.data[0].user].player.y = 32;
                    }
                });
            },
            handleMove: function(msg) {
                Vue.nextTick(function(){
                    var id = msg.data[0].user;
                    var x = parseInt(msg.data[0].x);
                    var y = parseInt(msg.data[0].y);
                    var z = parseInt(msg.data[0].z);
                    if(x != 0)
                    {
                        app.users[id].player.body.velocity.x = x;
                    }
                    if(y != 0)
                    {
                        app.users[id].player.body.velocity.y = y;
                    }
                    if(z != 0)
                    {
                        app.users[id].jumpTimer = app.game.time.now + z;
                    }
                });
            },
            handleMapLeave: function(msg) {
                Vue.nextTick(function(){
                    delete app.users[msg.data[0].user];
                });
            },
            enterMap: function () {
                this.submitCommand({
                    method: 'enter_map',
                    data:{}
                });
            },
            move: function (left,top,jump) {
                this.submitCommand({
                    method: 'move',
                    data:{
                        x: left.toString(),
                        y: top.toString(),
                        z: jump.toString(),
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