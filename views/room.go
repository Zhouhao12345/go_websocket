package views

import (
	"net/http"
	"go_ws/models"
	"go_ws/tools"
	"encoding/json"
	"log"
	"strings"
	"strconv"
)

func APIRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	signed, userId := tools.SingleSign(r)
	key := r.FormValue("key")
	if signed == false {
		http.Error(w, "Please sign in firstly!", http.StatusOK)
		return
	}

	//todo fixme improve it
	roomRows, err := models.SelectQuery(
		"select room.id as rid , k.user_names as des, k.images as images, room.create_date as room_create_date , " +
			"count(message.id) as unread, n.max_create_date from web_chatroom as room " +
			"inner join web_chatroom_users as chroom on room.id = chroom.chatroom_id " +
				"left join (select Max(message.create_date) as max_create_date, message.room_id from web_chatmessage as message " +
					"where message.user_id = ? group by message.room_id) as n on n.room_id = chroom.chatroom_id " +
					"left join web_chatmessage as message on message.room_id = chroom.chatroom_id and " +
					"message.unread = 1 and message.user_id != ?" +
						" inner join (SELECT GROUP_CONCAT(users.username) as user_names, " +
							"GROUP_CONCAT(gusers.avatar_image_small) AS images, " +
							"chroom.chatroom_id FROM auth_user AS users " +
								"INNER JOIN web_chatroom_users AS chroom ON chroom.ggacuser_id = users.id " +
									"INNER JOIN web_ggacuser AS gusers ON gusers.user_ptr_id = users.id " +
									"WHERE users.id != ? GROUP BY chroom.chatroom_id) as k on " +
										"k.chatroom_id = chroom.chatroom_id " +
						" where chroom.ggacuser_id = ? group by room.id having des like ? order by max_create_date desc",userId, userId, userId, userId, "%"+key+"%")
	for index, room := range roomRows {
		messageRows, err := models.SelectQuery(
			"select message.content, message.create_date from web_chatmessage as message " +
				"where message.room_id = ? order by create_date desc limit 1", room["rid"])
		if err != nil {
			log.Printf("error: %v", err)
		}

		if len(messageRows) > 0 {
			roomRows[index]["create_date"] = messageRows[0]["create_date"]
			roomRows[index]["content"] = tools.RemoveAllHtmlTags(messageRows[0]["content"])
		} else {
			roomRows[index]["create_date"] = room["room_create_date"]
			roomRows[index]["content"] = ""
		}
	}
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "DB ERROR", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.Encode(tools.ApiJsonNormalization(roomRows, 0, "success"))
	return
}

func APIRoomCreate(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	signed, userId := tools.SingleSign(r)
	if signed == false {
		http.Error(w, "Please sign in firstly!", http.StatusOK)
		return
	}

	r.ParseForm()
	userIds := r.Form["user_ids[]"]
	userIds = append(userIds, userId)
	userIdsSort, err := tools.Sort(userIds)
	userIdsStr := strings.Join(userIdsSort, "_")

	roomRaws, err := models.SelectQuery("select room.id as rid from web_chatroom as room where slug = ?",
		userIdsStr)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, "DB ERROR", http.StatusInternalServerError)
		return
	}

	if len(roomRaws) > 0 {
		json.NewEncoder(w).Encode(tools.ApiJsonNormalization(roomRaws, 0, "success"))
	} else {
		current_time := tools.Now().Format("2006-01-02 15:04:05")
		id ,err1 := models.InsertQuery("insert into web_chatroom VALUES (?,?,?,?,?,?,?)",
				0,userId, current_time, userId, current_time, "",  userIdsStr)
		if err1 != nil {
			log.Printf("error: %v", err1)
			http.Error(w, "DB ERROR", http.StatusInternalServerError)
			return
		}
		for _,user := range userIds {
			_ ,err2 := models.InsertQuery("insert into web_chatroom_users (chatroom_id, ggacuser_id) " +
				"VALUES (?,?)", id, user)
			if err2 != nil {
				log.Printf("error: %v", err2)
				http.Error(w, "DB ERROR", http.StatusInternalServerError)
				return
			}
		}

		rawRaws := make([]map[string]string, 0)
		rawRaws = append(rawRaws, map[string]string{
			"rid": strconv.FormatInt(id, 10),
		})
		json.NewEncoder(w).Encode(tools.ApiJsonNormalization(rawRaws, 0, "success"))
	}
	w.Header().Set("Content-Type", "application/json")
	return
}