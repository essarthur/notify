// https://gocodecloud.com/blog/2016/03/26/writing-a-text-adventure-game-in-go---part-2/


type Character struct {
	Name   string
	Health int
	Evasion int
	Alive  bool
	Speed  int
	Weap   int
	Npc    bool
}

func (p *Character) Equip(w int) {
	p.Weap = w
}

func (p *Character) Attack() int {
	return Weaps[p.Weap].Fire()
}

var ennemies = map[int]*Character{
	1: {Name: "Klingon", Health: 50, Alive: true, Weap: 2},
	2: {Name: "Romulan", Health: 55, Alive: true, Weap: 3},
}

type Weapon struct {
	minAtt int
	maxAtt int
	Name   string
}

func (w *Weapon) Fire() int {
	return w.minAtt + rand.Intn(w.maxAtt - w.minAtt)
}

var Weaps = map[int]*Weapon{
	1: {Name: "Phaser", minAtt: 5, maxAtt: 15},
	2: {Name: "Klingon Disruptor", minAtt: 1, maxAtt: 15},
	3: {Name: "Romulan Disruptor", minAtt: 3, maxAtt: 12},
}

type Players []Character

func (slice Players) Len() int {
	return len(slice)
}

func (slice Players) Less(i, j int) bool {
	return slice[i].Speed > slice[j].Speed //Sort descending
	//return slice[i].Speed < slice[j].Speed;		//Sort ascending
}

func (slice Players) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}


type Item struct {
	Name       string
	Action     string
	ItemForUse int
	Contains   []int
}

var Items = map[int]*Item{
	1: {Name: "Key"},
	2: {Name: "Chest", ItemForUse: 1, Contains: []int{3}},
	3: {Name: "Medal"},
}


func GetUserStrInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n >>> ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
	return text
}


//To be refactored on a location struct
func describeItems(player Character) {
	l := locMap[player.CurrentLocation]

	DisplayInfo("You see:")
	for _, itm := range l.Items {
		DisplayInfof("\t%s\n", Items[itm].Name)
	}
}


func ProcessCommands(player Character, input string) string {
	tokens := strings.Fields(input)
	command := tokens[0]
	itemName := ""
	if len(tokens) > 1 {
		itemName = tokens[1]
	}
	DisplayInfo(tokens)
	loc := locMap[player.CurrentLocation]
	switch command {
	case "get":
		//Make sure we do not pick it up twice
		if ItemInRoom(loc, itemName) && !ItemOnPlayer(player, itemName) {
			PutItemInPlayerBag(player, itemName)
			ItemRemoveFromRoom(loc, itemName)
		} else {
			DisplayInfo("Could not get " + itemName)
		}
	case "open":
		OpenItem(player, itemName)
	case "inv":
		DisplayInfo("Your Inventory: ")
		for _, itm := range player.Items {
			DisplayInfo("\t" + Items[itm].Name)
		}
	default:
	}
	return command
}

func main() {
	player = *new(Character)
	player.CurrentLocation = "main"
	input := ""
	for input != "quit" {
		describeItems(player)
		input = GetUserStrInput()
		input = ProcessCommands(player, input)
	}
}


func OpenItem(pla Character, itemName string) {
	DisplayInfo("Opening " + itemName)
	loc := locMap[player.CurrentLocation]
	for _, itm := range loc.Items {
		if Items[itm].Name == itemName {
			DisplayInfo("Loop >> " + Items[itm].Name)
			if Items[itm].ItemForUse != 0 && PlayerHasItem(pla, Items[itm].ItemForUse) {
				loc.Items = append(loc.Items, Items[itm].Contains...)
				Items[itm].Contains = *new([]int)
			}
		} else {
			DisplayInfo("Could not open the " + itemName)
		}
	}
}

func PlayerHasItem(pla Character, itm int) bool {
	for _, pitm := range pla.Items {
		if pitm == itm {
			return true
		}
	}
	return false
}

//To be refactored on a character struct
func PutItemInPlayerBag(pla Character, itemName string) {
	for index, itm := range Items {
		if itm.Name == itemName {
			player.Items = append(player.Items, index)
			return
		}
	}
}

//To be refactored on a location struct
func ItemRemoveFromRoom(loc *Location, itemName string) {
	for index, itm := range loc.Items {
		if Items[itm].Name == itemName {
			loc.Items = append(loc.Items[:index], loc.Items[index+1:]...)
		}
	}
}

//To be refactored on a location struct
func ItemInRoom(loc *Location, itemName string) bool {
	for _, itm := range loc.Items {
		if Items[itm].Name == itemName {
			return true
		}
	}
	return false
}

//To be refactored on a character struct
func ItemOnPlayer(pla Character, itemName string) bool {
	for _, itm := range pla.Items {
		if Items[itm].Name == itemName {
			return true
		}
	}
	return false
}


func RunBattle(players Players) {
	round := 1
	numAlive := players.Len()
	for {
		DisplayInfo("Combat round", round, "begins...")
		if endBattle(players) {
			break
		} else {
			DisplayInfo(players)
			round++
		}
	}
}

func endBattle(players []Character) bool {
	count := make([]int, 2)
	count[0] = 0
	count[1] = 0
	for _, pla := range players {
		if pla.Alive {
			if pla.Npc == false {
				count[0]++
			} else {
				count[1]++
			}
		}
	}
	if count[0] == 0 || count[1] == 0 {
		return true
	} else {
		return false
	}
}



func selectTarget(players []Character, x int) int {
	y := x
	for {
		y = y + 1
		if y >= len(players) {
			y = 0
		}
		if (players[y].Npc != players[x].Npc) && players[y].Alive {
			return y
		}
		if y == x {
			return -1
		}
	}
	return -1
}
With all that in place let see the new combat loop

func RunBattle(players Players) {
    sort.Sort(players)

	round := 1
	numAlive := players.Len()
	for {
		DisplayInfo("Combat round", round, "begins...")
        for x := 0; x < players.Len(); x++ {
            if players[x].Alive != true {
                continue
            }
            tgt := selectTarget(players, x)
            if tgt != -1 {
                DisplayInfo("player: ", x, "target: ", tgt)
                attp1 := players[x].Attack()
                players[tgt].Health = players[tgt].Health - attp1
                if players[tgt].Health <= 0 {
                    players[tgt].Alive = false
                    numAlive--
                }
                DisplayInfo(players[x].Name+" attacks and does", attp1, "points of damage with his", Weaps[players[x].Weap].Name, "to the ennemy.")
            }
        }
		if endBattle(players) {
			break
		} else {
			DisplayInfo(players)
			round++
		}
	}
}



func RunBattle(players Players) {
    sort.Sort(players)

	round := 1
	numAlive := players.Len()
	playerAction := 0
	for {
	    for x := 0; x < players.Len(); x++ {
    		players[x].Evasion = 0      // Reset evasion for all characters
    	}
		DisplayInfo("Combat round", round, "begins...")
        for x := 0; x < players.Len(); x++ {
            if players[x].Alive != true {
                continue
            }
            playerAction = 0
            if !players[x].Npc {
                DisplayInfo("DO you want to")
                DisplayInfo("\t1 - Run")
                DisplayInfo("\t2 - Evade")
                DisplayInfo("\t3 - Attack")
                GetUserInput(&playerAction)
            }
            if playerAction == 2 {
                players[x].Evasion = rand.Intn(15)
                DisplayInfo("Evasion set to:", players[x].Evasion)
            }
            tgt := selectTarget(players, x)
            if tgt != -1 {
                DisplayInfo("player: ", x, "target: ", tgt)
                attp1 := players[x].Attack()
                players[tgt].Health = players[tgt].Health - attp1
                if players[tgt].Health <= 0 {
                    players[tgt].Alive = false
                    numAlive--
                }
                DisplayInfo(players[x].Name+" attacks and does", attp1, "points of damage with his", Weaps[players[x].Weap].Name, "to the ennemy.")
            }
        }
		if endBattle(players) || playerAction == 1 {
			break
		} else {
			DisplayInfo(players)
			round++
		}
	}
}

func DisplayInfof(format string, args ...interface{}) {
	fmt.Fprintf(Out, format, args...)
}

func DisplayInfo(args ...interface{}) {
	fmt.Fprintln(Out, args...)
}

func GetUserInput(i *int) {
	fmt.Fscan(In, i)
}


