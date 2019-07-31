//javascript for h1.html
window.onload=function(){loadImg();};

//var Chess = new Array();
var Chess = [
	[0,0,0,0,0,0,0,0,0,0],
	[0,0,0,0,1,0,0,0,1,0],
	[0,0,0,0,0,1,0,0,0,1],
	[0,0,2,0,0,0,1,0,0,1],
	[0,0,2,0,0,0,0,0,0,0],
	[0,0,0,0,0,0,0,0,0,0],
	[0,2,0,0,0,0,1,0,1,0],
	[0,2,0,0,0,0,0,0,0,0],
	[0,0,0,1,0,0,0,0,0,0],
	[0,0,0,0,0,0,2,2,2,0],
];
var _Geo = {
	0: "wild",
	1: "river",
	2: "mountain",
	3: "forest",
};
var Animal = [
	[0,1,0,0,0,0,0,0,0,0],
	[0,0,0,0,0,0,0,0,2,0],
	[0,0,0,0,0,0,0,0,0,0],
	[0,0,0,0,0,0,0,0,3,0],
	[0,0,0,0,0,0,0,0,0,0],
	[0,0,0,0,0,0,0,0,0,0],
	[0,0,0,4,0,0,0,0,0,0],
	[0,0,0,0,0,0,0,0,0,0],
	[0,0,0,0,0,0,0,0,0,0],
	[0,0,0,0,0,0,0,0,0,0],
];   
var _Animals = {
	1: "deer",
	2: "frog",
	3: "tiger",
	4: "eagle",
	5: "",
	6: "",
	7: "",
	8: "",
	9: "",
	10: "",
	//used to draw elected role(half-transparent)
	101: "deer",   
	102: "frog",
	103: "tiger",
	104: "eagle",
	105: "",
	106: "",
	107: "",
	108: "",
	109: "",
	110: "",
};

//roles in chess
var _Roles = [];

//chess: space, animal, other
var _chess = {
	space: 0, 
	animal: "",
	other: "",
}

var _role = {
	name: "",
	country: 0,
	maxhp: 1,
	hp: 1,
	atk: 0,
	def: 0,
	cri: 0,
	spd: 0,
	st: 1,   //status:1-good, 2-heavy hurt, 0-dead
}

var _sidebarOpt = {
	opt: "",
	x: 0,
	y: 0,
	level: 1,   //level 0=opt unavailable, 1=normal, 2=this opt has several sub-opts
	preOpt: -1,   //record previous opt if level=2
}

var imgDeer = new Image();
imgDeer.src = "deer.png";
var imgFrog = new Image();
imgFrog.src = "frog.png";
var imgTiger = new Image();
imgTiger.src = "tiger.png";
var imgOwl = new Image();
imgOwl.src = "owl.png";
//
function loadImg(){
	imgDeer.addEventListener("load", function(){}, false);
	imgFrog.addEventListener("load", function(){}, false);
	imgTiger.addEventListener("load", function(){}, false);
	imgOwl.addEventListener("load", function(){}, false);
};

function loadRole(){
	for (var i=0;i<Animal.length;i++){
		var col = [];
		for (var j=0;j<Animal[0].length;j++){
			var role = Object.create(_role);
			if (Animal[i][j]!=0){
				role.name = _Animals[Animal[i][j]];
				switch (Animal[i][j]){
				case 1:
					role.country = 1;
					role.maxhp = 6;
					role.hp = role.maxhp;
					role.atk = 3;
					role.def = 3;
					role.spd = 3;
					break;
				case 2:
					role.country = 2;
					role.maxhp = 10;
					role.hp = role.maxhp;
					role.atk = 2;
					role.def = 2;
					role.spd = 3;
					break;
				case 3:
					role.country = 1;
					role.maxhp = 8;
					role.hp = role.maxhp;
					role.atk = 5;
					role.def = 5;
					role.spd = 4;
					break;
				case 4:
					role.country = 2;
					role.maxhp = 5;
					role.hp = role.maxhp;
					role.atk = 4;
					role.def = 4;
					role.spd = 5;
					break;
				default:
					break;
				}
			}
			col[j] = role;
		}
		_Roles[i] = col;
	}
}

//update role
function updateRole(p, prep){
	console.log("updateRole",p,prep);
	_Roles[parseInt(p/10)][p%10] = _Roles[parseInt(prep/10)][prep%10];
	_Roles[parseInt(prep/10)][prep%10] = Object.create(_role);
	return;
}

function p2xy(p){
	//var y = parseInt(p/10);
	//var x = p%10;
	return [p%10,parseInt(p/10)];
}

function xy2p(x,y){
	return y*10+x;
}

var c2 = document.querySelector("#cv2");
var ctx2 = c2.getContext("2d");
var prePos = 0;   //record previous position, p(x=p%10, y = parseInt(p/10))
c2.addEventListener("keyup", moveBlock, false);   //listen key up event
c2.addEventListener("mousedown", function(e){console.log("e");}, false);
c2.focus();   //listen mouse down event
var activeRole = 0;   //use to select a role, selected role
var objRoleP = -1;   //object role's position
var sidebar = false;   //if using sidebar
var sidebarMenu = new Array(4);   //_sidebarOpt array
var subMenu = [];
var preMenu = 0;   //previous option
var moveRange = -1;
var rangeFields = [];   //record P which in available range 
var atkRange = 1;

//canvas 2
function InitCV2_old(){
	turn = false;
	ctx2.clearRect(0,0,646,646);
	ctx2.fillStyle="#8836a0";
	for (var y=0;y<10;y++){
		for (var x=0;x<10;x++){
			if (turn){
				turn = !turn;
				ctx2.fillStyle="#8076a3";
			}else{
				turn = !turn;
				ctx2.fillStyle="#f9f1db";
			}
			ctx2.fillRect(x*64+3,y*64+3,64,64);
		}
		turn = !turn;
	}
};

function InitCV2(){
	//clear canvas2
	ctx2.clearRect(0,0,700,660);
	
	//draw background block
	for (var y=0;y<10;y++){
		for (var x=0;x<10;x++){
			switch (Chess[y][x]){   //different color with different space
			case 0:
				ctx2.fillStyle="#9abeaf";
				break;
			case 1:
				ctx2.fillStyle="#baccd9";
				break;
			case 2:
				ctx2.fillStyle="#b6a476";
				break;
			default:
				break;
			};
			ctx2.fillRect(x*64+3,y*64+3,64,64);
		};	
	};
	
	//draw geography line
	ctx2.strokeStyle="#702420";
	ctx2.lineWidth=0.5;
	for (var i=0;i<11;i++){   //h lines
		ctx2.beginPath();
		ctx2.moveTo(3,i*64+3);
		ctx2.lineTo(643, i*64+3);
		ctx2.stroke();
	};
	for (var j=0;j<11;j++){   //v lines
		ctx2.beginPath();
		ctx2.moveTo(j*64+3,3);
		ctx2.lineTo(j*64+3,643);
		ctx2.stroke();
	};
	InitAnimal();
};

//set animals on chess
function InitAnimal(){
	for (var y=0;y<10;y++){
		for (var x=0;x<10;x++){
			drawRole(Animal[y][x],x,y);
			if (Animal[y][x]!=0){   //if role is not empty
				switch (_Roles[y][x].country){   //draw country board
				case 1:
					ctx2.fillStyle = "#ce5777";
					ctx2.fillRect(x*64+3,(y+1)*64,64,4);
					break;
				case 2:
					ctx2.fillStyle = "#2486b9";
					ctx2.fillRect(x*64+3,(y+1)*64,64,4);
					break;
				default:
					break;
				}
			}
		};
	};
};

function drawRole(a,x,y){   //a is role type
	switch (a){
	case 1:
		ctx2.drawImage(imgDeer, x*64+3, y*64+3, 64, 64);
		break;
	case 2:
		ctx2.drawImage(imgFrog, x*64+3, y*64+3, 64, 64);
		break;
	case 3:
		ctx2.drawImage(imgTiger, x*64+3, y*64+3, 64, 64);
		break;
	case 4:
		ctx2.drawImage(imgOwl, x*64+3, y*64+3, 64, 64);
		break;
	default:
		break;
	};
}

//use wasd move
function moveBlock(e){
	if (activeRole!=0){   //selected animal -1 moverange when move
		var p = prePos;
		var y = parseInt(p/10);
		var x = p%10;
		// if (moveRange==0){
			// console.log("can't move now.")
			// return;
		// }else{
			// moveRange--;
		// }
		moveRange--;
		if (Animal[y][x]!=0&&activeRole==Animal[y][x]){   //clear previous block
			Animal[y][x]=0;
		}
	}
	switch (e.keyCode){
		case 38:   //up
			if (sidebar){
				moveUpSidebar();
			}else{
				moveRole("up");
			}
			break;
		case 87:   //w
			if (sidebar){
				moveUpSidebar();
			}else{
				moveRole("up");
			}
			break;
		case 40:   //down
			if (sidebar){
				moveDownSidebar();
			}else{
				moveRole("down");
			}
			break;
		case 83:   //s
			if (sidebar){
				moveDownSidebar();
			}else{
				moveRole("down");
			}
			break;
		case 37:   //left
			if (sidebar){
				moveLeftSidebar();
			}else{
				moveRole("left");
			}
			break;
		case 65:   //a
			if (sidebar){
				moveLeftSidebar();
			}else{
				moveRole("left");
			}
			break;
		case 39:   //right
			if (sidebar){
				moveRightSidebar();
			}else{
				moveRole("right");
			}
			break;
		case 68:   //d
			if (sidebar){
				moveRightSidebar();
			}else{
				moveRole("right");
			}
			break;
		case 13:   //enter
			if (sidebar){
				selectSidebar();   //enter to select option
			}else if (objRoleP!=-1){   //if an action affect two roles
				handleActions(objRoleP, prePos);
			}else{
				selectBlock();   //enter to select it
			}
			break;
		case 16:    //shift
			showStatus();   //shift to show status
			break;
		case 27:   //esc
			if (sidebar){   //if sidebar active, deactivate it, and select block
				sidebar = false;   //deactivate
				drawActiveBlock(prePos);   //draw block active
				selectBlock();   //select block
			}
			break;
		default:
			break;
	}
}

//draw active block which is selected
function drawActiveBlock(p){
	if (p < 0||p > 100){
		return;
	}
	//InitCV2();
	var y = parseInt(p/10);
	var x = p%10;
	//ctx2.strokeStyle="#2775b6";
	
	//ctx2.shadowBlur=20;
	//ctx2.shadowColor="#f2e7e5";
	
	if (activeRole!=0){
		if (Animal[y][x]==0){   //if select animal and space is empty
			Animal[y][x] = activeRole;
			//InitCV2();
			//InitAnimal();
		}else{
			drawRole(activeRole,x,y);
		}
	}
	//refresh chess and draw active board
	InitCV2();
	ctx2.strokeStyle="#fbf2e3";
	ctx2.lineWidth=3;
	ctx2.strokeRect(x*64+3,y*64+3,64,64)   //draw active board rect
	
	showInfo(x,y);
	console.log(_Roles[y][x]);
	//ctx2.drawImage(imgDeer, x*64+3, y*64+3, 64, 64);
	return;
};

//move Role
function moveRole(direction){
	var p = 0;
	switch (direction){
	case "up":
		p = prePos - 10;
		if (p < 0){
			return;
		}
		break;
	case "down":
		p = prePos + 10;
		if (p > 100){
			return;
		}
		break;
	case "left":
		p = prePos - 1;
		if (p%10==9||p < 0){
			return;
		}
		break;
	case "right":
		p = prePos + 1;
		if (p%10==0||p > 100){
			return;
		}
		break;
	default:
		break;
	}
	if (objRoleP!=-1){
		objRoleP = p;   //change objRoleP from origin role to new selected role
		drawActiveBlock(p);
		initARange();
	}else if (isRoleBlock(p)){   //can't select role block normally
		return;
	}else if (!canStay(p)){   //can't move to can't stay block
		return;
	}else if(activeRole!=0&&rangeFields.indexOf(p)==-1){   //role must move in valid ranges
		return;
	}else{
		if (activeRole!=0){
			updateRole(p, prePos);
		}
		prePos = p;
		drawActiveBlock(p);
		drawRange(p);
	};
};

function moveUpRole(){
	var p = prePos - 10;
	if (p < 0){
		return;
	}else if (isRoleBlock(p)){
		return;
	}else if (!canStay(p)){
		return;
	}else if(activeRole!=0&&rangeFields.indexOf(p)==-1){
		return;
	}else{
		prePos = p;
		drawActiveBlock(p);
		drawRange(p);
	};
};

function moveDownRole(){
	var p = prePos + 10;
	if (p > 100){
		return;
	}else if (isRoleBlock(p)){
		return;
	}else if (!canStay(p)){
		return;
	}else if(activeRole!=0&&rangeFields.indexOf(p)==-1){
		return;
	}else{
		prePos = p;
		drawActiveBlock(p);
		drawRange(p);
	};
};

function moveLeftRole(){
	var p = prePos - 1;
	if (p%10==9||p < 0){
		return;
	}else if (isRoleBlock(p)){
		return;
	}else if (!canStay(p)){
		return;
	}else if(activeRole!=0&&rangeFields.indexOf(p)==-1){
		return;
	}else{
		prePos = p;
		drawActiveBlock(p);
		// if (drawRange!=[]){
			// drawRange(p);
		// }
		drawRange(p);
	};
};

function moveRightRole(){
	var p = prePos + 1;
	if (p%10==0||p > 100){
		return;
	}else if (isRoleBlock(p)){
		return;
	}else if (!canStay(p)){
		return;
	}else if(activeRole!=0&&rangeFields.indexOf(p)==-1){
		return;
	}else{
		prePos = p;
		drawActiveBlock(p);
		drawRange(p);
	};
};

//init move range fields with position
function initRange(x,y){
	moveRange = _Roles[y][x].spd;
	rangeFields = [];
	for (var i=0;i<10;i++){
		for (var j=0;j<10;j++){
			if (Math.abs(i-x)+Math.abs(j-y)<=moveRange){
				rangeFields.push(j*10+i);
			}
		}
	}
	//console.log(x,y,rangeFields);
	drawRange(y*10+x);
}

//draw move track
function drawRange(p){
	for (var i=0;i<rangeFields.length;i++){
		//local block and role block not in range
		if (rangeFields[i]!=p&&!isRoleBlock(rangeFields[i])){
			var y = parseInt(rangeFields[i]/10);
			var x = rangeFields[i]%10;
			ctx2.beginPath();
			ctx2.strokeStyle = "#fbf2e3";
			ctx2.arc(x*64+35,y*64+35,24,0,2*Math.PI);
			ctx2.stroke();
		}
	}
}

//exchange the activeRole into new position
function exchangeRole2(preP, p){
	var ry = parseInt(preP/10);
	var rx = preP%10;
	var ey = parseInt(p/10);
	var ex = p%10;
	_Roles[ey][ex] = _Roles[ry][rx];
	_Roles[ry][rx] = Object.create(_role);
}

//check if it's empty block or role on block
function isRoleBlock(p){
	if (activeRole==0){
		return false;
	}
	var y = parseInt(p/10);
	var x = p%10;
	if (Animal[y][x]!=0){
		var otherinfo = "meet " + _Animals[Animal[y][x]] + ", battle or try another way";
		drawInfo(false, otherinfo);
		return true;
	}
	return false;
}

//active menu opt
function activeMenu(m){
	if (m<0||m>3){
		return;
	}
	
	ctx2.font = "16px Arial";
	for (var i=0;i<sidebarMenu.length;i++){   //draw highlight first, then draw menu
		if (i==m){
			ctx2.fillStyle = "#f9d27d";
			ctx2.fillRect(sidebarMenu[i].x-4,sidebarMenu[i].y-16,50,20);
		}else{
			ctx2.fillStyle = "#f1f0ed";
			ctx2.fillRect(sidebarMenu[i].x-4,sidebarMenu[i].y-16,50,20);
		}
		
		ctx2.fillStyle = "#503e2a";
		ctx2.fillText(sidebarMenu[i].opt, sidebarMenu[i].x, sidebarMenu[i].y);
	}

	if (sidebarMenu[m].preOpt>=0){
		//draw sub-menu
	}
}

function moveUpSidebar(){
	var m = preMenu-1;
	if (m<0){
		return;
	}else{
		preMenu = m;
		activeMenu(m);
	}
}

function moveDownSidebar(){
	var m = preMenu+1;
	if (m>3){
		return;
	}else{
		preMenu = m;
		activeMenu(m);
	}
}

function moveLeftSidebar(){
	var opt = sidebarMenu[preMenu];
	if (opt.level<2){
		return;
	}else if (opt.preOpt<0){
		return;
	}else{
		sidebarMenu[preMenu].preOpt=-1;
		activeMenu(preMenu);
	}
}

function moveRightSidebar(){
	var opt = sidebarMenu[preMenu];
	if (opt.level<2){
		return;
	}else if (opt.preOpt>=0){
		return;
	}else{
		sidebarMenu[preMenu].preOpt=preMenu;
		activeMenu(preMenu);
	}
}

function resourceReady(){
	return;
}

function putChess0(){
	for (var j=0;j<10;j++){
		var col = new Array();
		for (var i=0;i<10;i++){
			chess = Object.create(_chess);
			col[i] = chess;
		}
		Chess[j] = col;
	}
};

function putChess(){
	loadRole();
};

//deal with block select
function selectBlock(){
	var otherinfo = "";
	var p = prePos;
	var y = parseInt(p/10);
	var x = p%10;
	if (activeRole!=0){   //if select one role
		if (Animal[y][x]!=0){   //if block is not empty, battle
			otherinfo = _Animals[activeRole] + " vs " + _Animals[Animal[y][x]];
			drawInfo(false, otherinfo);
			//var result = initBattle();
			console.log(otherinfo);
		}else{
			if (canStay(p)){   //if role can stay in block
				Animal[y][x] = activeRole;
				activeRole = 0;
				InitCV2();
				if (!sidebar){
					rangeFields = [];
					initSidebar(x,y);   //if select a role and role in empty block
				}
			}else{   //else role can't stay in block
				otherinfo = _Animals[activeRole] + " can't stay in " + _Geo[Chess[y][x]];
				drawInfo(false, otherinfo);
				console.log(otherinfo);
			}
		}
	}else{
		if (Animal[y][x]!=0){   //if selected is empty and role in block, select it
			if (!sidebar){
				activeRole = Animal[y][x];
				Animal[y][x] = 0;
				moveRange = _Roles[y][x].spd;
				initRange(x,y);
			}
		}else{   
			//if selected and block both empty, show block info
		}
	}
};


//sidebar init
function initSidebar(x,y){
	var opt0 = Object.create(_sidebarOpt);
	opt0.opt = "Atk";
	opt0.x = x*64+72;
	opt0.y = y*64+16;
	sidebarMenu[0] = opt0;
	
	var opt1 = Object.create(_sidebarOpt);
	opt1.opt = "Skill";
	opt1.x = x*64+72;
	opt1.y = y*64+36;
	sidebarMenu[1] = opt1;
	
	var opt2 = Object.create(_sidebarOpt);
	opt2.opt = "Def";
	opt2.x = x*64+72;
	opt2.y = y*64+56;
	sidebarMenu[2] = opt2;
	
	var opt3 = Object.create(_sidebarOpt);
	opt3.opt = "Stay";
	opt3.x = x*64+72;
	opt3.y = y*64+76;
	sidebarMenu[3] = opt3;

	drawSidebarBG(x,y);
	activeMenu(0);
	//drawSidebarMenu();
	
	sidebar = true;
}

//draw sidebar background beside selected role
function drawSidebarBG(x,y){
	//draw country board and menu bg
	ctx2.fillStyle = "#f1f0ed";
	switch (_Roles[y][x].country){   
	case 1:
		ctx2.strokeStyle = "#ce5777";
		break;
	case 2:
		ctx2.strokeStyle = "#2486b9";
		break;
	default:
		ctx2.strokeStyle = "#b7ae8f";
		break;
	}
	ctx2.strokeRect(x*64+117,y*64+1,2,79);
	ctx2.strokeRect(x*64+68,y*64+79,50,2);
	ctx2.fillRect(x*64+69,y*64,49,80);
}

//deal with sidebar select
function selectSidebar(){
	var xy = p2xy(prePos);
	var role = _Roles[xy[1],xy[0]];
	drawActiveBlock(prePos);
	selectBlock();
	
	switch (preMenu){
	case 0:
		drawARange(prePos);   //draw
		break;
	case 1:
		role.atk *= 1.5;
		break;
	case 2:
		role.def *= 2;
		break;
	case 3:
		//do nothing
		break;
	default:
		break;
	}
	preMenu = 0;   //reset preMenu after select a option
	sidebar = false;
}

//init attack range
function drawARange(p){
	var xy = p2xy(p);
	objRoleP = p;
	rangeFields = [];
	for (var i=0;i<10;i++){
		for (var j=0;j<10;j++){
			if (Math.abs(i-xy[0])+Math.abs(j-xy[1])==1){   //if in attack range, push into rangeFields
				rangeFields.push(j*10+i);
			}
		}
	}
	initARange();
	return;
}


//draw attack range with red
function initARange(){
	for (var i=0;i<rangeFields.length;i++){
		//local block and role block not in range
		var y = parseInt(rangeFields[i]/10);
		var x = rangeFields[i]%10;
		ctx2.beginPath();
		ctx2.strokeStyle = "#ea7293";
		ctx2.arc(x*64+35,y*64+35,24,0,2*Math.PI);
		ctx2.stroke();
	}
	return;
}

//check if role can stay in space
function canStay(p){
	if (activeRole==0){
		return true;
	}
	var y = parseInt(p/10);
	var x = p%10;
	if (Chess[y][x]==1&&activeRole!=2){   //river
		return false;
	}
	if (Chess[y][x]==2&&activeRole!=4){   //mountain
		return false;
	}
	return true;
}

//show info on block
function showInfo(x,y){
	var mapinfo = "here is " + _Geo[Chess[y][x]];
	drawInfo(false, mapinfo);
	console.log(mapinfo);
	if (Animal[y][x]!=0){
		showStatus(x,y);
	}
}

//show role's status
function showStatus(x,y){
	var role = _Roles[y][x];
	console.log(role);
	var roleinfo = _Animals[Animal[y][x]]+" Atk:"+role.atk.toString()+" Def:"+role.def.toString()+" hp:"+role.hp.toString();
	drawInfo(false, roleinfo);
	//console.log(roleinfo);
}

//
function handleActions(objRoleP, userRoleP){
	if (preMenu==0){
		battle(objRoleP, userRoleP);
	}else{
		//deal with special action
	}
}

//battle with attacker and defender
function battle(defRoleP, atkRoleP){
	var result = initBattle(defRoleP, atkRoleP);
	var resInfo = "";
	if (result[2]){   //if success, defRole die, remove it from Animal, but still in _Roles
		resInfo = result[0].name + "die";
		drawInfo(false,resInfo);
		console.log(resInfo);
		Animal[parseInt(defRoleP/10)][defRoleP%10] = 0;
		//if beat, update role result
	}else{
		if (result[1].hp<=0){   //if atkRole die, remove it from Animal, but still in _Roles
			resInfo = result[1].name + "die";
			drawInfo(false,resInfo);
			console.log(resInfo);
			Animal[parseInt(atkRoleP/10)][atkRoleP%10] = 0;
		}
	}
	objRoleP = -1;   //reset object role position
	rangeFields = [];   //reset range fields array
	InitCV2();
}

//init battle between def role and atk role, if def role die, return success
function initBattle(defRoleP, atkRoleP){
	var defRole = _Roles[parseInt(defRoleP/10)][defRoleP%10];
	var atkRole = _Roles[parseInt(atkRoleP/10)][atkRoleP%10];
	var success = true;
	if (atkRole.atk-defRole.def>defRole.hp){
		defRole.hp = 0;
		defRole.st = 0;   //def role die
		return [defRole, atkRole, success]
	}else{
		defRole.hp = defRole.hp + defRole.def - atkRole.atk;
		if (defRole.maxhp/5>defRole.hp){   //if left hp < 20% maxhp, role heavy hurt
			defRole.st = 2; 
		}
		success = false;
	}
	if (defRole.spd>atkRole.spd){   //if def role speed > atk role, 
		if (defRole.atk-atkRole.def>atkRole.hp){
			atkRole.hp = 0;
			atkRole.st = 0;   //atk role die
		}else{
			atkRole.hp = atkRole.hp + atkRole.def - defRole.atk;
			if (atkRole.maxhp/5>atkRole.hp){
				atkRole.st = 2;   //atk role heavy hurt
			}
		}
	}
	
	return [defRole, atkRole, success]
}



//canvas 1
var c1 = document.querySelector("#cv1");
var ctx1 = c1.getContext("2d");
var icount = 0;
function InitCV1(){
	ctx1.fillStyle = "#fbf2e3";
	ctx1.strokeStyle = "#fbb957";
	ctx1.fillRect(0,0,320,660);
	ctx1.font = "60px Arial";
	ctx1.strokeText("Status Info",16,60);
	ctx1.strokeRect(16,72,290,1);
	//drawInfo("mapinfo","roleinfo");
}

//use ...arguments to accept any args
function drawInfo(refresh, ...args){
	var len = args.length;
	ctx1.fillStyle = "#fbb957";
	ctx1.font = "24px Arial";
	if (refresh){
		InitCV1();
		icount = 0;
	}
	for (var i=0;i<len;i++){   //fill text in each line, if height >maxheight, reset position to top
		if (120+40*icount>600){
			InitCV1();
			icount = 0;
		}
		ctx1.fillText(args[i], 16, 120+40*icount);
		icount++;
	}
}


//