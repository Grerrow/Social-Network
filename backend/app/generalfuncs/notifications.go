package generalfuncs

import (
	"log"
	"strconv"
	"time"

	"social-network/app/handlers/websocket"
	"social-network/app/models"
	"social-network/db"
)

func CreateNotification(n models.Notification) error {

	//horizontal expansion WORK BETTER
	_, err := db.Database.Exec(`INSERT INTO notifications 
	(user_id, follow_request_id, type, sender_id, sender_name, sender_avatar, post_id, group_id, group_name, event_id, is_read, created_at) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		n.UserID,
		n.FollowRequestID,
		n.Type,
		n.SenderID,
		n.SenderName,
		n.SenderAvatar,
		n.PostID,
		n.GroupID,
		n.GroupName,
		n.EventID,
		false,
		time.Now(),
	)

	//check if user is online and send websocket notification
	log.Printf("Creating notification for user %d of type %s", n.UserID, n.Type)
	if err != nil {
		log.Printf("Error creating notification: %v", err)
		return err
	}
	if websocket.IsUserOnline(strconv.Itoa(n.UserID)) {
		// for notification just a signal is enough
		websocket.SendToUser(
			strconv.Itoa(n.UserID),
			websocket.WebSocketMessage{
				Type: "notification",
				Data: "new_notification",
			},
		)

	}

	return err
}

// keep it simple for now and query the members of the group each time
// exclude the user who triggered the notification
// to avoid sending notification to self
// but i have to add the ws message for user that sent the message
func GetGroupMemberIDs(groupID int, excludeUserID int) ([]int, error) {
	var memberIDs []int
	rows, err := db.Database.Query(`SELECT user_id FROM group_members WHERE group_id = ? AND user_id != ?`, groupID, excludeUserID)
	if err != nil {
		log.Println("Error querying group members:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			log.Println("Error scanning group member ID:", err)
			continue
		}
		memberIDs = append(memberIDs, userID)
	}
	return memberIDs, nil
}
