//sync world info when websocket send back

package srv

import (
	""
)

type Position struct{
	P_level int
	P_block int   //blocks from start point to current point
	P_pos int   //real position
}

type World100 struct{
	Blocks int   //101 = 40 + 30 + 20 + 10 + boss
}

type World60 struct{
	Blocks int   //61 = 24 + 18 + 12 + 6 + boss
}

type World struct{
	LevelBlocks int
	TotalBlocks int
}

type SrvInfo struct{
	Players map[int]Player
	CurrentWorld World
}

var CurWorld World

func StartSrv(n int) (*SrvInfo, error){
	players := make(map[int]Player)
	world := World{LevelBlocks:n/10,TotalBlocks:n+1}
	srv := &SrvInfo{Players:players, CurrentWorld:world}
	return srv, nil
}
