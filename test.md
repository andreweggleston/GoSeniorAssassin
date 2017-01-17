It already returns:
    playerProfile:{
        "success": true,
        "data": {
            "id": 4,
            "studentid": "alexander_briasco",
            "name": "alexander briasco",
            "createdAt": "2016-11-17T12:57:38.389652-05:00",
            "tags": [
            "player"
            ],
            "role": "player",
            "bans": []
        }
    }

I want it to also return, inside of the "data" struct
  playerProfile:{
        "globalData":
        {
            "SafetyItem":
            "KillByDate":
            "Announcement":
        },
      "target": 'player',
      "markedForDeath": false,
      "killed": false
  }
and to return "target": in the regular data struct, that contains the "name" of that person's target8
and return the booleans "markedForDeath", and "killed". MfD should be set to true when the person who's
target is the current person calls assassinate()
and killed should be should be set to true when the player calls deathConfirmed()

Also need websocket strings for assassinate() and deathConfirmed()
