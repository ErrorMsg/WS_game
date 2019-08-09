var _player = {
	P_level:0,
	P_block:0,
	P_pos:0,   //position with level and block
	S_hp:10,   //add 10 hp each time move to higher level
	S_maxhp:10,
	S_atk:1,   //add 5 atk each time move to higher level
	S_def:1,   //add 5 def each time move to higher level
	S_spd:1,   //if defender spd fast then attacker, counter
	E_head:0,   //index of equipment
	E_weapon0:0,   //two-hands weapon 0
	E_weapon1:0,   //one-hand weapon 1
	E_weapon2:0,   //one-hand weapon 2
	E_body:0,
	E_shoes:0,
	E_jewelry1:0,   //player can dress two jewelries
	E_jewelry2:0,
	S_status:1,
	S_money:0,
	Handcards:[],
};

var _equip = {
	kind:0,
	id:0,
	hp:1,
	atk:1,
	def:1,
	spd:1,
	other:"",
}

var _monster = {
	id:0,
	hp:10,
	atk:3,
	def:3,
	spd:2,
	money:10,
}

var _boss = {
	id:0,
	hp:100,
	atk:30,
	def:30,
	spd:20,
}

var _experience = {
	kind:0,
	id:0,
}

var _bless = {
	id:0,
}

var _testplayer = {
	Name:"",
	Level:0,
	Block:1,
	Hands:[],
}

var c2 = document.querySelector("#cv2");
var ctx2 = c2.getContext("2d");
c2.addEventListener("keyup", move, false);
colors2 = ["#2983bb", "#1ba784", "#97846c", "#daa45a", "#f04b22"];   //blue, green, grey, orange, red
var colors = ["#baccd9", "#c6dfc8", "#f9f1db", "#f9e9cd", "#f2cac9"];
var activePlayer = null;

function init(){
	activePlayer = Object.create(_testplayer);
	activePlayer.Name = "echo";
	activePlayer.Level = 0;
	activePlayer.Block = 1;
	activePlayer.Hands = [];
	drawbg();
	draw0(40,40);
}

function drawbg(){
	for (var i=0;i<colors.length;i++){
		ctx2.fillStyle = colors[i];
		ctx2.fillRect(20+i*50,20,50,50);
	}
}

function draw0(x,y){
	ctx2.fillStyle = "#123456";
	ctx2.fillRect(x,y,10,10);
}

function move(e){
	switch (e.keyCode){
	case 13:   //enter
		break;
	case 68:   //d
		moveBlock(activePlayer);
	default:
		break;
	}
}

function moveBlock(activePlayer){
	activePlayer.Block ++;
	if (activePlayer.Block > colors.length){
		activePlayer.Level ++;
		activePlayer.Block = 1;
	}
	var x = activePlayer.Block * 50 - 10;
	var y = 40;
	drawbg();
	draw0(x,y);
	
	activePlayer.Hands[0] = activePlayer.Block*10;
	var msg = "j" + JSON.stringify(activePlayer);
	send(msg);
}

function Receive(e){
	console.log("Get: " + e.data);
}