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

type Timer struct {
    ThpTeleport int
    TnoMob int
    TuseItem int
    TuseSkill int
    TuseSkillSelf int
    TclickMove int
    TsameCoord int
    TsameMob int
    TsameItem int
    TclickLoot int
    TloadTP int
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
    CARTMIN = 0
    CARTMAX = 0

    SIT = false

    SSphere = 0


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

    mobDeadList = []Mob{}

    lockMap = ""
    saveMap = ""
    useTPNbAggro = 10
    useTPLockMap = -1
    useTPOnRoad = -1
    useTPDelay = 10
    useTPUnderHP = 5
    useSitUnderSP = -1
    useSitAboveSP = 99
    storageWeight = 49
    storageChoice = 2
    storageX = -1
    storageY = -1
    useGreed = -1
    useSphere = -1
    useSphereCombat = -1
    useHeal = -1
    useHealCombat = -1
    useHealLv = 10
    innSP = -1
    innX = -1
    innY = -1


    // ##### BOT
    lastMoveTime = time.Now().Unix()
    ccFrom = Coord{}
    ccTo = Coord{}
    pathTo = []Coord{ccFrom,ccTo}

    charCoord = Coord{}
    movePath = []Coord{}

    townRun = false
    innRun = false

    targetItemID = -1
    targetMobID = -1

    timers = Timer{
        ThpTeleport:0,
        TnoMob:0,
        TuseItem:0,
        TuseSkill:0,
        TuseSkillSelf:0,
        TclickMove:0,
        TsameCoord:0,
        TsameMob:0,
        TsameItem:0,
    }

    needWait = 0
    countAggro = 0

    sameCoord = Coord{}
    sameMob = 0
    sameItem = 0

)
