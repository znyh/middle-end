package config

import (
    "flag"
    "os"

    "kratos-demo/test/mygame/internal/base"

    "github.com/BurntSushi/toml"
    //"git.huoys.com/middle-end/kratos/pkg/log"
)

const GameID = 12001

var (
    ArenaID  = 1  //场ID: 1 2 3 4
    ServerID = "" //房间ID
)

var (
    //Ins 配置实例
    Ins = &config{}
)

func init() {
    flag.IntVar(&ArenaID, "aid", base.StrToInt(os.Getenv("ARENAID")), "specify the arena ID.")
    flag.StringVar(&ServerID, "sid", os.Getenv("HOSTNAME"), "specify the server ID.")
}

type (
    config struct {
        RC *RoomConfig
        CC *CardsDebugConfig
    }
    RoomConfig struct {
        GameCfg      GameConfig  `toml:"Game"`  //游戏配置
        TableCfg     TableConfig `toml:"Table"` //桌子配置
        RobotCfg     RobotConfig `toml:"Robot"` //机器人配置
        VipScoreRate float64
    }
    TableConfig struct {
        TableNum int `toml:"TableNum"` //桌子个数
        ChairNum int `toml:"ChairNum"` //椅子个数
    }
    //GameConfig 游戏配置
    GameConfig struct {
        DeBug    bool `toml:"DeBug"`    //debug模式
        LocalLog bool `toml:"LocalLog"` //本地日志

        MinMoney      int64   `toml:"MinMoney"`
        MaxMoney      int64   `toml:"MaxMoney"`
        EmojiCostList []int32 `toml:"EmojiCostList"` //6个道具的价格

        CardTimeout    int32 `toml:"CardTimeout"`    //基础时间  每轮重置
        HostingTimeout int32 `toml:"HostingTimeout"` //托管超时时间
        LaunchTimeout  int32 `toml:"LaunchTimeout"`  //加倍超时时间(s)
        ResultTime     int32 `toml:"ResultTime"`     //结算超时时间

        Marquee MarqueeConfig  `toml:"Marquee"`
        Coin    GameCoinConfig `toml:"Coin"`
    }
    GameCoinConfig struct {
        BaseScore       int64   `toml:"BaseScore"`       //底分
        TaiFlee         int64   `toml:"TaiFlee"`         //台费
        TaiFleeRate     float64 `toml:"TaiFleeRate"`     //台费 税率 默认值%10
        AvoidHurt       float32 `toml:"AvoidHurt"`       //免伤
        WinScoreLimit   int32   `toml:"WinScoreLimit"`   //赢分限定 整局游戏以双方中其中一方达到12分就结束游戏
        Multiple        int32   `toml:"Multiple"`        //游戏倍数最大可以增加到5倍 超过5倍的都按照5倍来计算。（可配置）
        SpecialScore    int32   `toml:"SpecialScore"`    //特殊的11分
        LaunchScoreList []int32 `toml:"LaunchScoreList"` //加倍列表
    }

    MarqueeConfig struct {
        WinMoney int64  `toml:"WinMoney"` //跑马灯单局的最小输赢
        Content  string `toml:"Content"`  //跑马灯消息提示 （祝贺 玩家名字 在 truco101 上赢得了 500,000 筹码！）
    }

    //配牌
    CardsDebugConfig struct {
        Enable bool    `toml:"Enable"` //是否开启
        Banker int32   `toml:"Banker"` //庄家椅子 [0,1,2,3]
        Hand1  []int32 `toml:"Hand1"`  //玩家1
        Hand2  []int32 `toml:"Hand2"`  //玩家2
        Hand3  []int32 `toml:"Hand3"`  //玩家3
        Hand4  []int32 `toml:"Hand4"`  //玩家4
        Next   []int32 `toml:"Next"`   //发牌列表
    }

    RobotConfig struct {
        IsOpen         bool             `toml:"IsOpen"`         //是否开机机器人
        TimeRange      []RobotTimeRange `toml:"RobotTimeRange"` //时间段配置
        MinMoney       int64            `toml:"MinMoney"`
        MaxMoney       int64            `toml:"MaxMoney"`
        TableMaxRobot  int32            `toml:"TableMaxRobot"`
        EnterRate      int32            `toml:"EnterRate"`
        StartEnterTime int32            `toml:"StartEnterTime"`
        PressInfo      RobotPressConfig `toml:"Press"`          //压测信息
        ThinkTimeRange []int32          `toml:"ThinkTimeRange"` //机器人思考时间区间
        LevelRange     []int32          `toml:"LevelRange"`     //加载机器人等级区间配置

    }

    RobotPressConfig struct {
        IsOpen bool `toml:"IsOpen"` //机器人压测
    }

    RobotTimeRange struct {
        Start string `toml:"Start"`
        End   string `toml:"End"`
        Num   int32  `toml:"Num"`
    }
)

func (c *RoomConfig) Set(txt string) error {
    roomcfg := &RoomConfig{}

    if err := toml.Unmarshal([]byte(txt), &roomcfg); err != nil {
        panic(err)
    }

    Ins.RC = roomcfg

    return nil
}

//LoadConfig 游戏配置
func LoadConfig() {
    //if err := paladin.Watch(fmt.Sprintf("gameconfig_%d.txt", ArenaID), Ins.RC); err != nil {
    //    panic(err)
    //}

    loadCardsCfg()

    return
}

//加载配牌文件
func loadCardsCfg() {

    if gameCfg := GetGC(); gameCfg != nil && gameCfg.DeBug == false {
        Ins.CC = &CardsDebugConfig{}
        return
    }

    cardsCfg := &CardsDebugConfig{}
    //if err := paladin.Get("cardsCfg.txt").UnmarshalTOML(&cardsCfg); err != nil {
    //    Ins.CC = &CardsDebugConfig{}
    //    return
    //}
    Ins.CC = cardsCfg
}

func GetIns() *config {
    return Ins
}

func GetGC() *GameConfig {
    return &Ins.RC.GameCfg
}

func GetTC() *TableConfig {
    return &Ins.RC.TableCfg
}

func GetRbC() *RobotConfig {
    return &Ins.RC.RobotCfg
}

func GetCC() *CardsDebugConfig {
    return Ins.CC
}

func GetStageInter(newStage int32) int64 {
    inter := int64(0)

    switch newStage {
    case StPrepare:
        inter = 0
    case StSendCard:
        inter = int64(5)
    case StOutCard:
        inter = int64(GetGC().CardTimeout)
    case StLaunch:
        inter = int64(GetGC().LaunchTimeout)
    case StResult:
        inter = int64(GetGC().ResultTime)
    default:
        inter = 15
    }

    return base.GetTick() + inter*1000
}
