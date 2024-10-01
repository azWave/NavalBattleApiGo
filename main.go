package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)


type Ship struct {
	Id int `json:"id"`
	Length int `json:"length"`
	PositionA [2]int `json:"positionA"`
	PositionB [2]int `json:"positionB"`
	Hits int `json:"hits"`
	Placed bool `json:"placed"`
}

type player struct {
	Id int `json:"id"`
	Board [10][10]int `json:"board"`
	Ships [5]Ship `json:"ships"`
}

type game struct {
	Id int `json:"id"`
	Players [2]player `json:"players"`
	
}

var ShipLengths = [5]int{5, 4, 3, 3, 2}
var games = []game{}

//function to create a new game and add it to the games list
func createGame(ctx *gin.Context) {
	var g game;
	g.Id = len(games) + 1
	for i := 0; i < 2; i++ {
		var p player
		p.Id = i + 1
		for j := 0; j < 5; j++ {
			var s Ship
			s.Id = j + 1
			s.Length = ShipLengths[j]
			s.Hits = 0
			s.Placed = false
			p.Ships[j] = s
		}
		g.Players[i] = p
	}
	games = append(games, g)

	//create jwt tokens for each player
	var tokens [2]string
	for i := 0; i < 2; i++ {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"gameId": g.Id,
			"playerId": i + 1,
		})

		tokenString, err := token.SignedString([]byte("secret"))
		if err != nil {
			ctx.IndentedJSON(500, gin.H{"error": "Failed to create token"})
			return
		}
		tokens[i] = tokenString
	}

	
	ctx.IndentedJSON(200, gin.H{"gameId": g.Id, "player1Token": tokens[0], "player2Token": tokens[1]})

}


func listGames(ctx *gin.Context) {
	ctx.IndentedJSON(200, games)
}


func gameInfo(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.IndentedJSON(401, gin.H{"error": "No token provided"})
		return
	}
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		ctx.IndentedJSON(401, gin.H{"error": "Invalid token"})
		return
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		ctx.IndentedJSON(401, gin.H{"error": "Invalid token"})
		return
	}
	gameId := int(claims["gameId"].(float64))
	playerId := int(claims["playerId"].(float64))
	if gameId < 1 || gameId > len(games) {
		ctx.IndentedJSON(400, gin.H{"error": "Invalid gameId"})
		return
	}
	if playerId < 1 || playerId > 2 {
		ctx.IndentedJSON(400, gin.H{"error": "Invalid playerId"})
		return
	}
	game := games[gameId - 1]
	player := game.Players[playerId - 1]
	var response = gin.H{"gameId": gameId, "playerId": playerId, "board": player.Board, "ships": player.Ships, "otherPlayerBoard": game.Players[1 - playerId % 2].Board}
	ctx.IndentedJSON(200, response)
}


func placeShip(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.IndentedJSON(401, gin.H{"error": "No token provided"})
		return
	}
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		ctx.IndentedJSON(401, gin.H{"error": "Invalid token"})
		return
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		ctx.IndentedJSON(401, gin.H{"error": "Invalid token"})
		return
	}
	gameId := int(claims["gameId"].(float64))
	playerId := int(claims["playerId"].(float64))
	if gameId < 1 || gameId > len(games) {
		ctx.IndentedJSON(400, gin.H{"error": "Invalid gameId"})
		return
	}
	if playerId < 1 || playerId > 2 {
		ctx.IndentedJSON(400, gin.H{"error": "Invalid playerId"})
		return
	}
	game := games[gameId - 1]
	player := game.Players[playerId - 1]

	PosA := ctx.Query("positionA")
	direction := ctx.Query("direction")
	ShipId := ctx.Query("shipId")
	if PosA == "" || PosB == "" || ShipId == "" {
		ctx.IndentedJSON(400, gin.H{"error": "Missing parameters"})
		return
	}
	//parse positionA
	var positionA [2]int
	_, err = fmt.Sscanf(PosA, "%d,%d", &positionA[0], &positionA[1])
	if err != nil {
		ctx.IndentedJSON(400, gin.H{"error": "Invalid positionA"})
		return
	}
	//parse direction

	_, err = fmt.Sscanf(direction, "%s", &direction)
	if err != nil {
		ctx.IndentedJSON(400, gin.H{"error": "Invalid direction"})
		return
	}
	CheckDirection := [4]string{"up", "down", "left", "right"}
	validDirection := false
	for _, dir := range CheckDirection {
		if direction == dir {
			validDirection = true
			break
		}
	}
	if !validDirection {
		ctx.IndentedJSON(400, gin.H{"error": "Invalid direction"})
		return
	}da

}



// Main function
func main() {
	router := gin.Default()

	router.GET("/createGame", createGame)
	router.GET("/gameInfo", gameInfo)
	router.GET("/Games", listGames)

	router.Run("localhost:8080")

}
