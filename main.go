package main

import (
	"io"
	"log"
	"os"
	"path"

	"github.com/real-mielofon/chess"
)

func GetTag(m *chess.Move) string {
	if m.HasTag(chess.Capture) {
		return "x"
	} else {
		return " "
	}
}

type FuncMoveCheck = func(m *chess.Move, mEthelon ElementMove) bool

type ElementMove struct {
	s1        chess.Square
	s2        chess.Square
	promo     chess.PieceType
	capture   bool
	checkType FuncMoveCheck
}

func CheckMove(m *chess.Move, mEthelon ElementMove) bool {
	return (m.S1() == mEthelon.s1) &&
		(m.S2() == mEthelon.s2) &&
		(m.Promo() == mEthelon.promo) &&
		(!mEthelon.capture || (mEthelon.capture && m.HasTag(chess.Capture)))
}

func CheckGame(g *chess.Game) bool {
	arrCheck := []ElementMove{
		ElementMove{s1: chess.C2, s2: chess.C4, promo: chess.NoPieceType, capture: false, checkType: CheckMove},
		ElementMove{s1: chess.C4, s2: chess.C5, promo: chess.NoPieceType, capture: false, checkType: CheckMove},
		ElementMove{s1: chess.C5, s2: chess.C6, promo: chess.NoPieceType, capture: false, checkType: CheckMove},
		ElementMove{s1: chess.C6, s2: chess.D7, promo: chess.NoPieceType, capture: true, checkType: CheckMove},
	}

	current := 0
	for i, move := range g.Moves() {
		if i%2 == 0 {
			if CheckMove(move, arrCheck[current]) {
				if current > 2 {
					log.Printf("%d: %d\n", i, current)
				}
				current += 1
				if current > len(arrCheck)-1 {
					return true // дошли до конца
				}
			} else {
				if move.HasTag(chess.Capture) {
					p := g.Positions()[i].Board().Piece(move.S2())
					if p == chess.BlackQueen {
						log.Printf("%d: убита черная королева\n", i, current)
						return false // убита черная королева
					}
				}
			}
		} else if (move.S2() == arrCheck[current].s1) && move.HasTag(chess.Capture) {
			log.Printf("%d: наша фигура убита на шаге %d\n", i, current)
			return false // убита наша фигура
		}
	}
	return false
}

func CheckPNG(filePGN string, outPngFileName string, logFieName string) (count int, countCheckSuccess int) {
	f, _ := os.Open(filePGN)
	defer f.Close()
	log.Println("Process PGN = ", filePGN)

	logFile, err := os.OpenFile(logFieName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	pgnFile, err := os.OpenFile(outPngFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer pgnFile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	//	game := chess.NewGame()
	games, err := chess.GamesFromPGN(f)
	if err != nil {
		log.Println("GamesFromPGN with err = ", err)
	}
	count, countCheckSuccess = 0, 0
	for i, game := range games {
		count++
		if CheckGame(game) {
			countCheckSuccess++
			board := game.Position().Board().Draw()
			log.Printf("FIND!! %d: %d: \n\n", count, i)
			pgnFile.WriteString(game.String())
			pgnFile.WriteString("\n\n")
			log.Println(game.String())
			log.Println("")
			log.Println(board)
			log.Println("")
		}
	}
	log.Printf("PreSummary %d: %d:", count, countCheckSuccess)
	return count, countCheckSuccess
}

func ExamplePrintGMGames() {
	arrPgns := []string{
		"Databases/test.pgn",
		/*
			"Databases/ficsgamesdb_search_53480.pgn",
			"Databases/ficsgamesdb_search_53479.pgn",
			"Databases/ficsgamesdb_search_53478.pgn",
			"Databases/ficsgamesdb_search_53477.pgn",
			"Databases/ficsgamesdb_search_53475.pgn",
			"Databases/ficsgamesdb_search_53481.pgn",
		*/
	}

	count := 0
	countCheckSuccess := 0

	for _, filePGN := range arrPgns {
		localCount, localCountCheckSuccess := CheckPNG(filePGN, "pgnfind.pgn", "log.txt")
		count += localCount
		countCheckSuccess += localCountCheckSuccess
	}
	log.Printf("Summary %d: %d:", count, countCheckSuccess)
}

func main() {
	if len(os.Args) < 2 {
		ExamplePrintGMGames()
		return
	}
	pngFileName := os.Args[1]
	ext := path.Ext(pngFileName)

	outPngFileName := pngFileName[0:len(pngFileName)-len(ext)] + ".out.png"
	if len(os.Args) >= 3 {
		outPngFileName = os.Args[2]
	}

	logFileName := pngFileName[0:len(pngFileName)-len(ext)] + ".log"
	if len(os.Args) >= 4 {
		logFileName = os.Args[3]
	}

	count, countCheckSuccess := CheckPNG(pngFileName, outPngFileName, logFileName)
	log.Printf("Summary %d: %d:", count, countCheckSuccess)
}
