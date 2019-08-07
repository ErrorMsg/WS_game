//define actions
//html动作 选择起点，移动，买卖，加入手牌，弃牌，选择攻击对象攻击方法，
//srv动作 初始化地图，洗牌shuffle，初始化人物，翻牌pickcard，处理事件牌handlerexp，商店，战斗，休息，使用牌handleruse
//发送data typelen + [[jsontype1,datalen1],[jsontype2,datalen2]...]+jsondata1+jsondata2+...

package srv

import (
	"encoding/json"
)

//handle move action
type ActionMove struct{
	PID int   //Player.Link.L_ID
	Level int   //Player.Position.P_level
	Block int
}

func handlerMove(data []byte) error{
	var info ActionMove
	err := json.Unmarshal(data, &info)
	player := PlayerList[info.PID]
	player.P_block += 1
	if player.P_block  > (CurWorld.LevelBlocks*(4-player.P_level)){
		player.P_level += 1
		if player.P_level > 4{
			player.P_level = 4
		}
		player.P_block = 1
	}
	//send back json.Marshal(player)
	return nil
}


//handle pick card action
type ActionPick struct{
	PID int
	CardType string
	Count int
}

func handlerPick(data []byte) error{
	var info ActionPick
	err := json.Unmarshal(data, &info)
	switch info.CardType{
	case "monster":
		cs,err := moncards.pickCard(info.Count)
		if err != nil{
			return err
		}
		//send back pick cards []Card
	case "experience":
		cs, err := expcards.pickCard(info.Count)
		if err != nil{
			return err
		}
		//send back pick cards []Card
		if len(cs) == 1{
			handleExperience(info.PID, cs[0])
		}else{
			//wait player choose one exp, and feed back
		}
		
	case "equipment":
		cs, err := equipcards.pickCard(info.Count)
		if err != nil{
			return err
		}
		//send back pick cards []Card
	default:
		
	}
}

//most choose action on html
type ActionChoose struct{
	PID int
	CardID int
}

func handlerChoose(data []byte) error{
	var info ActionChoose
	err := json.Unmarshal(data, &info)
	//get card from db by info.CardID
	//card := 
	switch card.C_kind{   //card.Card.C_kind
	case 0:   //used with 0-expcard normally
		handleExperience(info.PID, card)
	case 1:   //1-moncard
		//init battle
	default:   //equipcard & freecard
		player := PlayerList[info.PID]
		player.HandCards = append(player.HandCards, card)
		PlayerList[info.PID] = player
		//send back playerlist
	}
}

func handleExperience(pid int, c Card) error{
	//handle expcard effect
	sc, err := c.SpecificCard()
	if err != nil{
		return err
	}
	switch sc.C_effect{
	case "battle":
		mc, err := moncards.pickCard(1)
		//send back mc
	case "shop":
		ec, err := equipcards.pickCard(3)
		fc, err := freecards.pickCard(2)
		//send back ec and fc
	case "rest":
		player := PlayerList[pid]
		player.S_hp += player.S_maxhp/4
		if player.S_hp > player.S_maxhp{
			player.S_hp = player.S_maxhp
		}
		//send back player list
	default:
	}
}

//handle battle with monster action
type ActionBattle struct{
	PID int
	MonsterID int
}

func handlerBattle(data []byte) error{
	var info ActionBattle
	err := json.Unmarshal(data, &info)
	if err != nil{
		return err
	}
}


//handle vs another player action
type ActionVS struct{
	AttackorID int
	DefenderID int
}

func handlerVS(data []byte) error{
	var info ActionVS
	err := json.Unmarshal(data, &info)
	if err != nil{
		return err
	}
}

//deal with use handcard action
type ActionUse struct{
	PID int
	CardID int
}

func handlerUse(data []byte) error{
	var info ActionUse
	err := json.Unmarshal(data, &info)
	if err != nil{
		return err
	}
	player := PlayerList[PID]
	//get card from db by info.CardID
	//card := 
	if card.C_kind == 2{   //card.Card.C_kind
		switch card.Part{   //1-head,2-two hand weapon,3-one hand weapon,4-body,5-shoes,6-jewelry
		case 1:
			player.E_head = card.C_id
		case 2:
			player.E_weapon0 = card.C_id
		case 3:
			if player.E_weapon1 != 0 && player.E_weapon2 != 0{
				
			}else if player.E_weapon1 != 0{
				player.E_weapon2 = card.C_id
			}else{
				player.E_weapon1 = card.C_id
			}
		case 4:
			player.E_body = card.C_id
		case 5:
			player.E_shoes = card.C_id
		case 6:
			if player.E_jewelry1 != 0 && player.E_jewelry2 != 0{
				
			}else if player.E_jewelry1 != 0{
				player.E_jewelry2 = card.C_id
			}else{
				player.E_jewelry1 = card.C_id
			}
		}
	}else if card.C_kind == 3{
	
	}else{
	
	}
}

//common player list update
type ActionUpdate struct{
	players []Player
}

func handlerUpdate(data []byte) error{
	var info ActionUpdate
	err := json.Unmarshal(data, &info)
	if err != nil{
		return err
	}
	for p := range info.players{
		pid := p.L_ID
		PlayerList[pid] = p
	}
}