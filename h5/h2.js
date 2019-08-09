var _player = {
	level:1,
	block:0,
	pos:0,   //position with level and block
	hp:10,   //add 10 hp each time move to higher level
	atk:1,   //add 5 atk each time move to higher level
	def:1,   //add 5 def each time move to higher level
	spd:1,   //if defender spd fast then attacker, counter
	eq_head:0,   //index of equipment
	eq_weapon0:0,   //two-hands weapon 0
	eq_weapon1:0,   //one-hand weapon 1
	eq_weapon2:0,   //one-hand weapon 2
	eq_body:0,
	eq_shoot:0,
	eq_jewelry1:0,   //player can dress two jewelries
	eq_jewelry2:0,
	st:1,
	money:0,
	handcards:[],
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

var Experiences = [];
var DExperiences = [];   //discard experience cards

var Equips = [];
var Equips0 = [];   //1st level
var Equips1 = [];   //2nd level
var Equips2 = [];   //3rd level
var Equips3 = [];   //4th level

//var DEquips = [];   //discard equip cards
var DEquips0 = [];
var DEquips1 = [];
var DEquips2 = [];
var DEquips3 = [];

var c2 = document.querySelector("#cv2");
var ctx2 = c2.getContext("2d");
c2.addEventListener("keyup", test, false);
colors2 = ["#2983bb", "#1ba784", "#97846c", "#daa45a", "#f04b22"];   //blue, green, grey, orange, red
colors = ["#baccd9", "#c6dfc8", "#f9f1db", "#f9e9cd", "#f2cac9"];
var activePlayer = null;


ctx2.fillStyle = "#123456";
ctx2.fillRect(60,60,50,50);
ctx2.fillStyle = "#fedcba";
ctx2.fillRect(50,50,30,30);

if (ctx2.isPointInPath(120,120)){
	console.log("true");
}

//draw 5 level background cycles in full map
function gen100(){   //full map
	for (var i=0;i<5;i++){
		var r = 300 - i*60;   //each cycle -60 radius
		var n = 40 - i*10;   //each cycle -10 blocks
		if (n<=1){
			n = 1;
		}
		ctx2.fillStyle = colors[i];
		ctx2.beginPath();
		ctx2.arc(300,300,r,0,2*Math.PI);
		ctx2.closePath();
	
		ctx2.fill();
		drawNCycle(300,300,r,n);
	}
}

//draw 5 level background cycles in middle map
function gen60(){   //middle map
	for (var i=0;i<5;i++){
		var r = 300 - i*60;
		var n = 24 - i*6;
		if (n<=1){
			n = 1;
		}
		ctx2.fillStyle = colors[i];
		ctx2.beginPath();
		ctx2.arc(300,300,r,0,2*Math.PI);
		ctx2.closePath();
	
		ctx2.fill();
		drawNCycle(300,300,r,n);
	}
	//drawStartPoint();
	drawPoint(p1.level, p1.block);
}

//draw sector boards
function drawNCycle(x,y,r,n){
	for (var i=0;i<n;i++){
		ctx2.beginPath();
		if (n!=1){
			ctx2.moveTo(x,y);
		}
		ctx2.arc(x,y,r,2*Math.PI/n*i,2*Math.PI/n*(i+1));
		ctx2.closePath();
		ctx2.stroke();
	}
}

//get [x,y], n = total block in level, m = current block, r = radius
function getXYC(n,m,r){   
	var x = 300 - (r-30)*Math.cos((m+0.5)/n*2*Math.PI);
	var y = 300 - (r-30)*Math.sin((m+0.5)/n*2*Math.PI);
	return [x,y];
}

function getXY(n,m,r){
	var x = 300 - (r-20)*Math.cos((m+0.2)/n*2*Math.PI);
	var y = 300 - (r-20)*Math.sin((m+0.2)/n*2*Math.PI);
	return [x,y];
}


//draw start points according to maps
function drawStartPoint(){
	var xy = [];
	ctx2.fillStyle = "#ffffff";
	for (var i=0;i<6;i++){   //max player = last level
		ctx2.beginPath();
		xy = getXY(24,i*4,300);
		ctx2.arc(xy[0],xy[1],10,0,2*Math.PI);
		ctx2.closePath();
		ctx2.fill();
	}	
}


function drawPoint(n,m){
	var xy = getXY(24-n*6,m,300-60*n);
	console.log(xy);
	ctx2.fillStyle = "#ffffff";
	ctx2.beginPath();
	ctx2.arc(xy[0],xy[1],10,0,2*Math.PI);
	ctx2.closePath();
	ctx2.fill();	
}

//init equip cards and experience cards when start
function initCards(){
	var r = 0;
	for (var i=0;i<Experiences.length;i++){
		r = random(i,Experiences.length);
		var temp = Experiences[i];
		Experiences[i] = Experiences[r];
		Experiences[r] = temp;
	}
	for (var j=0;j<Equips.length;j++){
		r = random(j,Equips.length);
		var temp = Equips[j];
		Equips[j] = Equips[r];
		Equips[r] = temp;
	}
}

function random(s,b){
	return Math.floor(Math.random()*(b-s))+s;
}

function swap(a,b){
	var temp = a;
	a = b;
	b = temp;
	return [a,b];
}

var p1 = Object.create(_player);
p1.block = 12;
activePlayer = p1;

function test(e){
	switch (e.keyCode){
	case 13:
		moveBlock(activePlayer);
		break;
	default:
		break;
	}
}

function moveBlock(p){
	p.block++;
	if (p.block>=24-p.level*6){
		p.block = 0;
		p.level++;
	}
	gen60();
	drawPoint(p.level, p.block);
	console.log("moving to "+p.block);
	if (p.block != 0){
		pickExp(p);
	}else if (p.level==5){
		vsBoss(p);
		p.level = 0;
	}else{
		initLevel(p.level);
		pickExp(p);
	}
	
}

function pickExp(p){
	var exp = Experiences.shift();   //get a card from top
	handleExp(p,exp);
	
}

function handleExp(p,exp){
	switch (exp.kind){
	case "monster":
		initMonster(p,exp);
		break;
	case "shop":
		initShop(p);
		break;
	case "something":
		initSomething(p,exp);
		break;
	}
}

function initShop(p){
	var es = [];
	switch (p.level){
	case 0:
		es = Equips0;
		break;
	case 1:
		es = Equips1;
		break;
	case 2:
		es = Equips2;
		break;
	case 3:
		es = Equips3;
		break;
	default:
		break
	}
	var e1 = es.shift();   //display first three equip cards
	var e2 = es.shift();
	var e3 = es.shift();
	displayCards(e1,e2,e3);
}


function useCard(p,card){
	if (p.handcards.indexof(card)==-1){
		return;
	}
	switch (card.kind){
	case "equip":
		break;
	case "bless":
		break;
	case "":
		break;
	default:
		break;
	}
	removeCard(p,card);
}

function removeCard(p,card){
	var i = p.handcards.indexof(card);
	//var cc = p.handcards[0:i].concat(p.handcards[i+1:p.handcards.length]);
	var cc = [];
	for (var j=0;j<i;j++){
		cc[j] = p.handcards[j];
	}
	for (var j=i;j<p.handcards.length-1;j++){
		cc[j]=p.handcards[j+1];
	}
	p.handcards = cc;
	drawHandcards(p);
}

function drawHandcards(p){
	//show player's handcards with img
}
	