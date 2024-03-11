package main

import(
    "sync"
    "time"
)

type Coord struct {
	X,Y int
}

type Mob struct {
    MobID int
    CoordsFrom Coord
    CoordsTo Coord
    PathMoveTo []Coord
    // HPMax int
    // HPLeft int
    LastMoveTime int64
    MoveSpeed int
    DeathTime int64
    IsNotValid bool
    Priority int
    Aggro bool
    IsLooter bool
    AtSight bool
    Name string
    Bexp int
    Jexp int
    TPdist int
}

type Npc struct {
    NpcID int
    Coords Coord
    Name string
}

type Item struct {
    ItemID int
    Coords Coord
    Amount int
    DropTime int64
    IsValid bool
    Priority int
    EqSlot int
}

type Player struct {
    Name string
    Coords Coord
}

type Trap struct {
    TrapID int
    Coords Coord
    Radius int
}

var(
    accountID = 0


    MOVESPEED = 0
    BASEXPMAX = 0
    BASEEXP = 0
    JOBXPMAX = 0
    JOBEXP = 0
    CHARNAME = ""
    BASELV = 0
    JOBLV = 0
    ZENY = 0
    MAP = ""
    HPLEFT = 0
    HPMAX = 0
    WEIGHTMAX = 0
    WEIGHT = 0
    SPLEFT = 0
    SPMAX = 0

    SIT = false



    MUnpcList sync.Mutex
    npcList = map[int]Npc{}
    MUmobList sync.Mutex
    mobList = map[int]Mob{}
    MUgroundItems sync.Mutex
    groundItems = map[int]Item{}
    MUinventoryItems sync.Mutex
    inventoryItems = map[int]Item{}
    MUstorageItems sync.Mutex
    storageItems = map[int]Item{}
    MUcartItems sync.Mutex
    cartItems = map[int]Item{}
    MUbuffList sync.Mutex
    buffList = map[int][]int64{}
    MUplayerList sync.Mutex
    playerList = map[int]Player{}
    MUtrapList sync.Mutex
    trapList = map[int]Trap{}

    MUmobDeadList sync.Mutex
    mobDeadList = map[int]Mob{}

    SSphere = 0


    lockMap = ""
    saveMap = ""
    useTPNbAggro = 10
    useTPLockMap = 0
    useTPOnRoad = 0
    useTPDelay = 10
    useTPUnderHP = 5
    useSitUnderSP = 0
    useSitAboveSP = 99

    storageWeight = 49

    storageX = 0
    storageY = 0
    useGreed = 0


    // ##### BOT
    lastMoveTime = time.Now().Unix()
    ccFrom = Coord{}
    ccTo = Coord{}
    pathTo = []Coord{}

    charCoord = Coord{}
    movePath = []Coord{}

    townRun = false

    targetItemID = -1
    targetMobID = -1

    needWait = 0

    noMobTimer = 0
    useItemTimer = 0
    clickMoveTimer = 0

    chkTimecharCoord = 0
    chkTimetargetMobID = 0
    chkTimetargetItemID = 0

    chkcharCoord = Coord{}
    chktargetMobID = 0
    chktargetItemID = 0

)
