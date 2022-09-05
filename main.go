package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// --------------------------------------------->Player<--------------------------------------------------START
type player struct {
	currentPosition string
	inventory       []string
}

//инициализируем структуру с игроком

var player1 player // создаем глобальную переменную где будем хранить структуру игрока, изначально все поля равны значению по умолчанию

func initPlayer() {
	player1.currentPosition = "кухня"
}

//функция для создания начальных паарметров игрока( по дефолту только позиция на кухне)

// --------------------------------------------->Player<-----------------------------------------------------END
// -------------------------------------------------------------------------------------------------------------
// --------------------------------------------->World<----------------------------------------------------START
// создаем мапу с команатами их соедржимым и связями
var world map[string]map[string][]string

// функция возвращающая мапу со всеми локациями и их содержимым
func initWord() map[string]map[string][]string {
	kitchen := make(map[string][]string)
	kitchen["connecting"] = []string{"коридор"}
	kitchen["itemTake"] = []string{"Нож", "Вилка"}

	hall := make(map[string][]string)
	hall["connecting"] = []string{"кухня", "комната", "улица"}

	myRoom := make(map[string][]string)
	myRoom["connecting"] = []string{"коридор"}
	myRoom["itemTake"] = []string{"ключи", "конспект"}
	myRoom["itemPutOn"] = []string{"рюкзак", "шапка"}

	world1 := make(map[string]map[string][]string)
	world1["кухня"] = kitchen
	world1["коридор"] = hall
	world1["комната"] = myRoom
	return world1
}

// создаем мапу с закрытими локациями
var closedLocation map[string]string

// создаем закрытые локации
func initClosedLocation() map[string]string {
	closedLocation := make(map[string]string)
	closedLocation["улица"] = "ключи"

	return closedLocation
}

//--------------------------------------------->World<-----------------------------------------------------END
//------------------------------------------------------------------------------------------------------------
//--------------------------------------------->ACTIVITY<-------------------------------------------------START

func proccessingRequest(command string) {
	//разбиваем строку на команды
	commandsWords := strings.Split(command, " ")
	switch len(commandsWords) {
	case 1:
		oneCommandAction(commandsWords)
	case 2:
		twoCommandAction(commandsWords)
	case 3:
		threeCommandAction(commandsWords)
	}
}

func oneCommandAction(command []string) {
	switch command[0] {
	case "оглядеться", "осмотреться":
		lookAround()
	case "инвентарь":
		fmt.Println("Ваш инвентарь: |", player1.inventory, "|")
	}

}

func twoCommandAction(command []string) {
	switch command[0] {
	case "идти":
		changeLocation(command[1])
	case "взять":
		take(command[1], "itemTake")
	case "надеть":
		take(command[1], "itemPutOn")
	}
}

func threeCommandAction(command []string) {
	switch command[0] {
	case "применить", "Применить":
		use(command)
	}
}

//-------------------------------------------->ACTIONS<----------------------------------------------------

func use(command []string) {
	//проверка на существование предмета в иневнтаре
	itemIndex := stringExistInSlice(player1.inventory, command[1])
	if itemIndex == -1 {
		fmt.Println("У меня нет  \"", command[1], "\"")
		return
	}
	//проверили на то что находимся рядом с объектом
	objInThisLoca := stringExistInSlice(world[player1.currentPosition]["connecting"], command[2])

	if objInThisLoca == -1 {
		fmt.Println("Не вижу", command[2], ".Это точно то место?")
		return
	}

	key := closedLocation[command[2]]
	//проверили на что объекты могут взаимодейтвовать
	if key == command[1] {
		delete(closedLocation, command[2]) //удаляем ключ закрытой двери из мапы
		player1.inventory = delObjFromSlice(player1.inventory, command[1])
		fmt.Println("Взаимодействие успешно!")
		return
	}
	fmt.Println("Взаимодействие не удалось")
} //ПРИМЕНИТЬ

func lookAround() {
	fmt.Println("|-----------------------------------------------------------------|")
	fmt.Print("Текущее местоположение : ", player1.currentPosition, ".")
	printMapProperty("itemTake", "Оглядевшись вокруг вы заметили : ")
	printMapProperty("itemPutOn", "Еще из интересного : ")
	printMapProperty("connecting", "Вы можете пойти по следующим направлениям: ")
	fmt.Println("|-----------------------------------------------------------------|")
} // ОСМОТРЕТЬСЯ

func take(itemName string, partLocation string) {
	objectExist := checkObjectInLocation(itemName, partLocation)
	// првоеряю есть ли этот предмет в текущей локации

	if objectExist {
		player1.inventory = append(player1.inventory, itemName)
		deleteObjectWithWorld(itemName, partLocation)
		fmt.Println("Добавлено в инвентарь :", itemName)
	} else {
		fmt.Println("Не удается поднять ", itemName)
	}
	// предмет есть - добавляю его в инвентарь и удаляю из комнаты
	// предмета нет - пишу что невозможо взять предмет
} // ВЗЯТЬ

func changeLocation(nextLocation string) {
	//сравниваем есть ли в текущей локации сообщение с тем куда мы идем
	connectingExist := checkObjectInLocation(nextLocation, "connecting")

	// если есть то перемещаемся (меняем текузую локацию в структуре игрока)
	if connectingExist && !thisIsClosedLocation(nextLocation) {
		player1.currentPosition = nextLocation
		lookAround()
	} else {
		fmt.Print("Невозможно пройти к ", nextLocation, ".")
		if thisIsClosedLocation(nextLocation) && connectingExist {
			fmt.Println("Нужно открыть дверь.")
		}
	}
	// если нет то пишем (невозможно переместиться)
} // ИДТИ

//--------------------------------------------->UTILS<------------------------------------------------------

func stringExistInSlice(sl []string, str string) int {
	for i := 0; i < len(sl); i++ {
		if sl[i] == str {
			return i
		}
	}
	return -1
}

func printMapProperty(propetyName string, textBefore string) {
	exist := len(world[player1.currentPosition][propetyName])
	if exist != 0 {
		fmt.Print(textBefore)
		for ; exist > 0; exist-- {
			fmt.Print(world[player1.currentPosition][propetyName][exist-1])
			if exist != 1 {
				fmt.Print(",")
			}
		}
		fmt.Print(".\n")
	} else if propetyName == "itemTake" {
		fmt.Print("Ничего интересного. \n")
	}
}

func checkObjectInLocation(object string, partLocation string) bool {
	amountLocation := len(world[player1.currentPosition][partLocation])

	//сравниваем есть ли в текущей локации сообщение с тем куда мы идем
	for ; amountLocation > 0; amountLocation-- {
		if object == world[player1.currentPosition][partLocation][amountLocation-1] {
			return true
		}
	}
	return false
}

func thisIsClosedLocation(locationName string) bool {
	_, ok := closedLocation[locationName]
	return ok
}

func delObjFromSlice(sl []string, str string) []string {
	indexItem := stringExistInSlice(sl, str)
	sl[indexItem] = sl[len(sl)-1]
	sl[len(sl)-1] = ""
	sl = sl[:len(sl)-1]
	return sl
}

func deleteObjectWithWorld(itemName string, partLocation string) {
	world[player1.currentPosition][partLocation] = delObjFromSlice(world[player1.currentPosition][partLocation], itemName)
}

//------------------------------------------>MAIN<-----------------------------------------------------------

func main() {

	initGame()
	command := bufio.NewScanner(os.Stdin)
	fmt.Println("|-----------------------------------------------------------------|")
	fmt.Println("|--------Приветсвую тебя! Советую тебе снчала осмотреться!--------|")
	fmt.Println("|--------Команды : идти,взять,инвентарь,надеть,применить----------|")
	fmt.Println("|--------Твоя цель выйти из дома!---------------------------------|")
	fmt.Println("|-----------------------------------------------------------------|")
	for {
		if player1.currentPosition == "улица" {
			fmt.Println("Ну теперь катим в универ")
			go func() {
				for {
					for _, r := range `-\|/` {
						fmt.Printf("\r%c", r)
						time.Sleep(time.Millisecond * 100)
					}
				}
			}()
		}
		command.Scan()
		line := command.Text()
		proccessingRequest(line)
	}
}

func initGame() {

	initPlayer()
	world = initWord()
	closedLocation = initClosedLocation()
	/*
		эта функция инициализирует игровой мир - все команты
		если что-то было - оно корректно перезатирается
	*/
}
