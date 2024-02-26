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
    Coords Coord
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
    AtSight bool
    Name string
}

type Item struct {
    ItemID int
    Coords Coord
    Amount int
    DropTime int64
    IsValid bool
    Priority int
}

var(
    accountID = 0

    BASEXPMAX = 0
    BASEEXP = 0
    JOBXPMAX = 0
    JOBEXP = 0
    CHARNAME = ""
    BASELV = 0
    JOBLV = 0
    ZENY = 0
    MAP = ""
    XPOS = 0
    YPOS = 0
    HPLEFT = 0
    HPMAX = 0
    WEIGHTMAX = 0
    WEIGHT = 0
    SPLEFT = 0
    SPMAX = 0

    SIT = false

    needWait = 0
    now = time.Now()

    MUmobList sync.Mutex
    mobList = map[int]Mob{}
    MUgroundItems sync.Mutex
    groundItems = map[int]Item{}
    MUinventoryItems sync.Mutex
    inventoryItems = map[int]Item{}
    MUbuffList sync.Mutex
    buffList = map[int][]int64{}

    SSphere = 0


    lockMap = ""
    saveMap = ""
    useTPNbAggro = 10
    useTPNbAggroLoot = 10
    useTPLockMap = 0
    useTPDelay = 10
    useSitUnderSP = 0
    useSitAboveSP = 99
    timerNoMob = 0

    // ##### BOT
    charCoord = Coord{}
    nextPoint = Coord{}
    nextStep = Coord{}
    movePath = []Coord{}
    pathIndex = 0
    minDist = 1
    distFromDest = float64(0)

    targetItemID = -1
    targetMobID = -1

    chkTimecharCoord = 0
    chkTimetargetMobID = 0
    chkTimetargetItemID = 0

    chkcharCoord = Coord{}
    chktargetMobID = 0
    chktargetItemID = 0

)
