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
    PathTo []Coord
    HPMax int
    HPLeft int
    LastMoveTime int64
    MoveSpeed int
    DeathTime int64
    IsNotValid bool
}

type Item struct {
    ItemID int
    Coords Coord
    Amount int
    DropTime int64
    IsValid bool
}

var(

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

    needWait = 0
    now = time.Now()
    stateTime = time.Now()

    MUmobList sync.Mutex
    mobList = map[int]Mob{}
    MUgroundItems sync.Mutex
    groundItems = map[int]Item{}
    MUinventoryItems sync.Mutex
    inventoryItems = map[int]Item{}
    MUbuffList sync.Mutex
    buffList = map[int][]int64{}

    SSphere = 0

    // lastMobDead = Mob{}
    // lastItemLooted = -1


    accountID = 0
    lockMap = ""
    saveMap = ""
    useTPLockMap = 0
    useTPDelay = 10

    // ##### BOT
    charCoord = Coord{}
    nextPoint = Coord{}
    nextStep = Coord{}
    movePath = []Coord{}
    pathIndex = 0
    minDist = 1

    // targetMob = -1
    // targetItem = -1
    //
    // attackDist = 2
    // attackIndex = 0

)
