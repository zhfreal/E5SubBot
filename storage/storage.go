package storage

import (
	"fmt"
	"time"

	"github.com/zhfreal/E5SubBot/config"
	"github.com/zhfreal/E5SubBot/logger"
	ms "github.com/zhfreal/E5SubBot/microsoft"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"

	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var (
		err  error
		dial gorm.Dialector
	)

	switch config.DB {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Mysql.User,
			config.Mysql.Password,
			config.Mysql.Host,
			config.Mysql.Port,
			config.Mysql.DB,
		)

		if config.Mysql.SSLMode != "" {
			dsn += "&tls=" + config.Mysql.SSLMode
		}

		dial = mysql.Open(dsn)
	case "sqlite":
		dial = sqlite.Open(config.Sqlite.DB)
	}

	if dial == nil {
		logger.Fatalln("failed to get dial, please check your config")
	}
	DB, err = gorm.Open(dial, &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now()
		},
		// NamingStrategy: schema.NamingStrategy{
		//     TablePrefix: config.Table_Prefix,
		// },
	})
	if err != nil {
		logger.Fatalf("failed to open db, failed with: %v\n", err.Error())
	}
	DB.AutoMigrate(&AppConfig{})
	DB.AutoMigrate(&Users{})
	DB.AutoMigrate(&Stats{})
	DB.AutoMigrate(&Op{})
	DB.AutoMigrate(&OpDetails{})
	DB.AutoMigrate(&TaskRecords{})
	// Init Ops
	InitOps()
	// init Stats
	InitStats()
}

const (
	Table_APPs        = "apps"
	Table_Users       = "users"
	Table_Stats       = "stats"
	Table_OpDetails   = "op_details"
	Table_TaskRecords = "task_records"
	Table_Op          = "op"
)

// ///////////////////////////////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////

type TypeUserIDOpID struct {
	UserId uint
	OpId   uint
}

// just for stat
type AppConfigStat struct {
	ID       uint
	ClientId string
	Alias    string
	Bound    int
}

// just for stat notification
type ResultStatsForNotify struct {
	AppAlias   string `gorm:"column:app_alias;not null"`
	UserAlias  string `gorm:"column:user_alias;not null"`
	OpAlias    string `gorm:"column:op_alias;not null"`
	TgId       int64  `gorm:"column:tg_id;not null"`
	Success    int    `gorm:"column:success;not null"`
	Failure    int    `gorm:"column:failure;not null"`
	LatestTime int64  `gorm:"column:last_time"`
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////
// AppConfig db helpers
// ///////////////////////////////////////////////////////////////////////////////////////////////////////
type AppConfig struct {
	gorm.Model
	ClientId     string `gorm:"column:client_id;not null"`
	ClientSecret string `gorm:"column:client_secret"`
	TenantID     string `gorm:"column:tenant_id"`
	Alias        string `gorm:"column:alias;not null"`
	FromLegacy   bool   `gorm:"column:from_legacy;type:boolean;default:false"`
	Enabled      bool   `gorm:"column:enabled;type:boolean;default:true"`
}

func (a *AppConfig) TableName() string {
	return Table_APPs
}

func AddAppConfig(a *AppConfig) error {
	return DB.Create(a).Error
}

func UpdateAppConfig(a *AppConfig) error {
	return DB.Save(a).Error
}

func DelAppConfig(id uint) error {
	return DB.Table(Table_APPs).Where("id = ?", id).Delete(&AppConfig{}).Error
}

func DelAppConfigByClientID(client_id string) error {
	return DB.Table(Table_APPs).Where("client_id = ?", client_id).Delete(&AppConfig{}).Error
}

func GetAllApp() ([]*AppConfig, error) {
	var configs []*AppConfig
	err := DB.Table(Table_APPs).Find(&configs).Error
	return configs, err
}

func GetAllAppEnabled() ([]*AppConfig, error) {
	var configs []*AppConfig
	err := DB.Table(Table_APPs).Find(&configs).Where("enabled = ?", 1).Error
	return configs, err
}

func GetAppByID(id uint) (*AppConfig, error) {
	var config AppConfig
	t := DB.Table(Table_APPs).Where("id = ?", id).First(&config)
	err := t.Error
	if t.RowsAffected == 0 {
		return nil, err
	}
	return &config, err
}

func GetAppIDByClientID(client_id string) (uint, error) {
	config, e := GetAppByClientID(client_id)
	if e != nil {
		return 0, e
	} else if config == nil {
		return 0, fmt.Errorf("can't get app ID by client_id %s", client_id)
	}
	return config.ID, nil
}

func GetAppByClientID(client_id string) (*AppConfig, error) {
	var config AppConfig
	t := DB.Table(Table_APPs).Where("client_id = ?", client_id).First(&config)
	err := t.Error
	if t.RowsAffected == 0 {
		return nil, err
	}
	return &config, err
}

func GetAppByAlias(alias string) (*AppConfig, error) {
	var config AppConfig
	t := DB.Table(Table_APPs).Where("alias = ?", alias).First(&config)
	err := t.Error
	if t.RowsAffected == 0 {
		return nil, err
	}
	return &config, err
}

// "id": app_id
func AppsHas(id uint) bool {
	return !(DB.Table(Table_APPs).
		Where("id = ?", id).
		First(&AppConfig{}).RowsAffected == 0)
}

func AppsHasClientID(client_id string) bool {
	return !(DB.Table(Table_APPs).
		Where("client_id = ?", client_id).
		First(&AppConfig{}).RowsAffected == 0)
}

func AppsHasAlias(alias string) bool {
	return !(DB.Table(Table_APPs).
		Where("alias = ?", alias).
		First(&AppConfig{}).RowsAffected == 0)
}

func AppsHasClientIDOrAlias(client_id, alias string) bool {
	return AppsHasClientID(client_id) || AppsHasAlias(alias)
}

// get specific app stats by it's ID, used by admin
func GetAppStatByID(id uint) *AppConfigStat {
	all_user_stat := GetAppsStat()
	for _, stat := range all_user_stat {
		if stat.ID == id {
			return stat
		}
	}
	return nil
}

// get specific app stats by it's alias, used by admin
func GetAppStatByAlias(alias string) *AppConfigStat {
	all_user_stat := GetAppsStat()
	for _, stat := range all_user_stat {
		if stat.Alias == alias {
			return stat
		}
	}
	return nil
}

// get all Apps stats, used by admin
func GetAppsStat() map[string]*AppConfigStat {
	var stats []*AppConfigStat
	DB.Table(Table_APPs).Select(Table_APPs + ".id, " + Table_APPs + ".client_id, " + Table_APPs + ".alias, count(*) as bound").
		Joins("right join " + Table_Users + " on " + Table_APPs + ".id = " + Table_Users + ".app_id").
		Group(Table_APPs + ".id, " + Table_APPs + ".client_id, " + Table_APPs + ".alias").Find(&stats)
	m := make(map[string]*AppConfigStat)
	for _, stat := range stats {
		m[stat.Alias] = stat
	}
	return m
}

// get all apps without any user bound, return []*AppConfigStat
func GetAppStatsOnlyNoUserBound() []*AppConfigStat {
	all_apps := GetAppsStat()
	var stats []*AppConfigStat
	for _, app := range all_apps {
		if app.Bound == 0 {
			stats = append(stats, app)
		}
	}
	return stats
}

// get all apps with users bound, return []*AppConfigStat
// used by admin only
func GetAppStatsOnlyBoundUser() []*AppConfigStat {
	all_apps := GetAppsStat()
	var stats []*AppConfigStat
	for _, app := range all_apps {
		if app.Bound > 0 {
			stats = append(stats, app)
		}
	}
	return stats
}

// get all apps with users bound, return []*AppConfigStat
func GetAppStatsForUser(tg_id int64) []*AppConfigStat {
	all_apps := GetAppsStatByTgID(tg_id)
	var stats []*AppConfigStat
	for _, app := range all_apps {
		if app.Bound > 0 {
			stats = append(stats, app)
		}
	}
	return stats
}

// get all Apps stats, these app already bound a account by the specific user's tg_id
func GetAppsStatByTgID(tg_id int64) map[string]*AppConfigStat {
	var stats []*AppConfigStat
	DB.Table(Table_APPs).Select(Table_APPs+".id, "+Table_APPs+".client_id, "+Table_APPs+".alias, count(*) as bound").
		Where(Table_Users+".tg_id = ?", tg_id).
		Joins("right join " + Table_Users + " on " + Table_APPs + ".id = " + Table_Users + ".app_id").
		Group(Table_APPs + ".id, " + Table_APPs + ".client_id, " + Table_APPs + ".alias").Find(&stats)
	m := make(map[string]*AppConfigStat)
	for _, stat := range stats {
		m[stat.Alias] = stat
	}
	return m
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////
// table users
// ///////////////////////////////////////////////////////////////////////////////////////////////////////
type Users struct {
	gorm.Model
	AppId        uint   `gorm:"column:app_id;not null"`
	TgId         int64  `gorm:"column:tg_id;not null"`
	MsId         string `gorm:"column:ms_id;not null"`
	MsUsername   string `gorm:"column:ms_username;not null"` // ms account, an email address
	Alias        string `gorm:"column:alias;not null"`       // display name of ms account
	AccessToken  string `gorm:"column:access_token;not null"`
	RefreshToken string `gorm:"column:refresh_token;not null"`
	ExpiresAt    int64  `gorm:"column:expire_at;not null"`
	LastTime     int64  `gorm:"column:last_time;not null"`
	Enabled      bool   `gorm:"column:enabled;type:boolean;default:true"`
}

func (u *Users) TableName() string {
	return Table_Users
}

func AddUsers(c *Users) error {
	return DB.Create(c).Error
}

func UpdateUsers(c *Users) error {
	return DB.Save(c).Error
}

func UpdateUsersTokens(id uint, uc *Users) error {
	c, e := GetUsersByID(id)
	if e != nil {
		return e
	}
	c.AccessToken = uc.AccessToken
	c.RefreshToken = uc.RefreshToken
	c.ExpiresAt = uc.ExpiresAt
	return DB.Save(c).Error
}

func DelUsers(id uint) error {
	return DB.Table(Table_Users).Where("id = ?", id).Delete(&Users{}).Error
}

func GetAllUsers() ([]*Users, error) {
	var configs []*Users
	err := DB.Table(Table_Users).Find(&configs).Error
	return configs, err
}

func GetAllUsersEnabled() ([]*Users, error) {
	var configs []*Users
	err := DB.Table(Table_Users).
		Joins("join "+Table_APPs+" on "+Table_APPs+".id = "+Table_Users+".app_id").
		Where(Table_APPs+".enabled = ? AND "+Table_Users+".enabled = ?", 1, 1).
		Find(&configs).Error
	return configs, err
}

func GetUsersByID(id uint) (*Users, error) {
	var config Users
	t := DB.Table(Table_Users).Where("id = ?", id).First(&config)
	err := t.Error
	if t.RowsAffected == 0 {
		return nil, err
	}
	return &config, err
}

func GetUsersByAppIDAndMsID(app_id, ms_id uint) (*Users, error) {
	var config Users
	t := DB.Table(Table_Users).Where("app_id = ? AND ms_id = ?", app_id, ms_id).First(&config)
	err := t.Error
	if t.RowsAffected == 0 {
		return nil, err
	}
	return &config, err
}

func GetUsersByAppID(app_id uint) ([]*Users, error) {
	var configs []*Users
	err := DB.Table(Table_Users).Where("app_id = ?", app_id).Find(&configs).Error
	return configs, err
}

func GetUsersByAppIDUserAlias(app_id uint, alias string) (*Users, error) {
	var config *Users
	t := DB.Table(Table_Users).Where("app_id = ? AND alias = ?", app_id, alias).Find(&config)
	err := t.Error
	if t.RowsAffected == 0 {
		return nil, err
	}
	return config, err
}

func GetUsersByClientID(client_id string) ([]*Users, error) {
	var configs []*Users
	t_query := fmt.Sprintf("select %s.* from %s, %s where %s.app_id = %s.id and %s.client_id = ?",
		Table_Users, Table_Users, Table_APPs, Table_Users, Table_APPs, Table_APPs)
	t_query = fmt.Sprintf("%s and %s.deleted_at IS NULL and %s.deleted_at IS NULL",
		t_query, Table_APPs, Table_Users)
	err := DB.Raw(t_query).Scan(&configs).Error
	return configs, err
}

func GetUsersByAppIDTgID(app_id uint, tg_id int64) ([]*Users, error) {
	var configs []*Users
	err := DB.Table(Table_Users).Where("app_id = ? AND tg_id = ?", app_id, tg_id).Find(&configs).Error
	return configs, err
}

func GetUserByAppIDTgIDMsUsername(app_id uint, tg_id int64, ms_username string) ([]*Users, error) {
	var config []*Users
	t := DB.Table(Table_Users).Where("app_id = ? AND tg_id = ? AND ms_username = ?", app_id, tg_id, ms_username).First(&config)
	if t.RowsAffected == 0 || t.Error != nil {
		return nil, t.Error
	}
	return config, nil
}

func GetUserByAppIDTgIDUserAlias(app_id uint, tg_id int64, user_alias string) (*Users, error) {
	var config *Users
	t := DB.Table(Table_Users).Where("app_id = ? AND tg_id = ? AND alias = ?", app_id, tg_id, user_alias).First(&config)
	if t.RowsAffected == 0 || t.Error != nil {
		return nil, t.Error
	}
	return config, nil
}

func GetUsersByTgID(tg_id int64) ([]*Users, error) {
	var configs []*Users
	err := DB.Table(Table_Users).Where("tg_id = ?", tg_id).Find(&configs).Error
	return configs, err
}

func HasUserByAppIDTgIDMsUsername(app_id uint, tg_id int64, ms_username string) bool {
	found, _ := FindUserByAppIDTgIDMsUsername(app_id, tg_id, ms_username)
	return found
}

func FindUserByAppIDTgIDMsUsername(app_id uint, tg_id int64, ms_username string) (bool, *Users) {
	var config Users
	t := DB.Table(Table_Users).Where("app_id = ? AND tg_id = ? AND ms_username = ?", app_id, tg_id, ms_username).
		First(&config).RowsAffected
	return t == 1, &config
}

func UpdateUserByAppIDTgIDMsUsername(uc *Users) error {
	app_id := uc.AppId
	tg_id := uc.TgId
	ms_username := uc.MsUsername
	new_uc := &Users{
		MsId:         uc.MsId,
		Alias:        uc.Alias,
		AccessToken:  uc.AccessToken,
		RefreshToken: uc.RefreshToken,
		ExpiresAt:    uc.ExpiresAt,
	}
	return DB.Table(Table_Users).Where("app_id = ? AND tg_id = ? AND ms_username = ?", app_id, tg_id, ms_username).Updates(&new_uc).Error
}

// get all UsersConfig by match with given email in UsersConfig.MsUsername
func GetUsersConfigIDByEmail(email string) ([]*Users, error) {
	var config []*Users
	err := DB.Table(Table_Users).Where("ms_username = ?", email).Find(&config).Error
	return config, err
}

// make new UsersConfig instance from app_id, tg_id, and cached_token
func NewUsersConfig(app_id uint, tg_id int64, cached_token *ms.TokenCache) *Users {
	return &Users{
		AppId:        app_id,
		TgId:         tg_id,
		MsId:         cached_token.ClientID,
		MsUsername:   cached_token.Username,
		Alias:        cached_token.Alias,
		AccessToken:  cached_token.AccessToken,
		RefreshToken: cached_token.RefreshToken,
		ExpiresAt:    cached_token.ExpireTime,
	}
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////
// table op
// ///////////////////////////////////////////////////////////////////////////////////////////////////////
type Op struct {
	gorm.Model
	Alias string `gorm:"column:alias;not null"`
}

func (o *Op) TableName() string {
	return Table_Op
}

func AddOp(o *Op) error {
	return DB.Create(o).Error
}

func AddOps(os []*Op) error {
	return DB.Create(os).Error
}

func DelOp(id uint) error {
	return DB.Table(Table_Op).Where("id = ?", id).Delete(&Op{}).Error
}

func UpdateOp(o *Op) error {
	return DB.Save(o).Error
}

func GetAllOps() ([]*Op, error) {
	Ops := make([]*Op, 0)
	err := DB.Table(Table_Op).Find(&Ops).Error
	return Ops, err
}

func GetOpByID(id uint) (*Op, error) {
	op := &Op{}
	t := DB.Table(Table_Op).Where("id = ?", id).First(&op)
	err := t.Error
	if t.RowsAffected == 0 {
		return nil, err
	}
	return op, err
}

func GetOpAliasByID(id uint) (string, error) {
	op, e := GetOpByID(id)
	if e != nil {
		return "", e
	}
	return op.Alias, nil
}

func InitOps() {
	// add config to Op
	ops, e := GetAllOps()
	if e != nil {
		logger.Fatalf("failed to get ops, failed with: %v\n", e.Error())
		panic(e)
	}
	// init op_dict
	op_dict := make(map[string]uint)
	for _, op := range ops {
		op_dict[op.Alias] = op.ID
	}
	// loop OPs
	op_slice := make([]*Op, 0)
	for id, op := range ms.OPs {
		// no-exists
		if _, ok := op_dict[op]; !ok {
			t_op := &Op{
				Model: gorm.Model{ID: id},
				Alias: op,
			}
			t_op.ID = id
			op_slice = append(op_slice, t_op)
		}
	}
	// no OPs to add
	if len(op_slice) == 0 {
		return
	}
	e = AddOps(op_slice)
	if e != nil {
		logger.Fatalf("failed to store ops, failed with: %v\n", e.Error())
		panic(e)
	}
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////
// table stats
// ///////////////////////////////////////////////////////////////////////////////////////////////////////
type Stats struct {
	gorm.Model
	UserID   uint  `gorm:"column:user_id;not null"`
	OpID     uint  `gorm:"column:op_id;not null"`
	Success  int   `gorm:"column:success;not null"`
	Failure  int   `gorm:"column:failure;not null"`
	LastTime int64 `gorm:"column:last_time;not null"` // in milliseconds
}

func (s *Stats) TableName() string {
	return Table_Stats
}

func AddStats(o *Stats) error {
	return DB.Create(o).Error
}

func AddStatsSlice(ops []*Stats) error {
	return DB.Create(ops).Error
}

func InitStats() error {
	ops, e := GetAllOps()
	if e != nil {
		return e
	}
	users, e := GetAllUsers()
	if e != nil {
		return e
	}
	opstats, e := GetAllStats()
	if e != nil {
		return e
	}
	opstats_dict := make(map[TypeUserIDOpID]bool, 0)
	for _, opstat := range opstats {
		opstats_dict[TypeUserIDOpID{UserId: opstat.UserID, OpId: opstat.OpID}] = true
	}
	this_time := time.Now().Unix()
	op_stat_slice := make([]*Stats, 0)
	for _, u := range users {
		for _, op := range ops {
			new_user_op := TypeUserIDOpID{UserId: u.ID, OpId: op.ID}
			if _, ok := opstats_dict[new_user_op]; !ok {
				op_stat_slice = append(op_stat_slice, &Stats{
					UserID:   u.ID,
					OpID:     op.ID,
					Success:  0,
					Failure:  0,
					LastTime: this_time,
				})
			}
		}
	}
	if len(op_stat_slice) > 0 {
		return DB.Create(op_stat_slice).Error
	}
	return nil
}

func InitStatsByUserID(uid uint) error {
	ops, e := GetAllOps()
	if e != nil {
		return e
	}
	this_time := time.Now().Unix()
	op_stat_slice := make([]*Stats, 0)
	for _, op := range ops {
		op_stat_slice = append(op_stat_slice, &Stats{
			UserID:   uid,
			OpID:     op.ID,
			Success:  0,
			Failure:  0,
			LastTime: this_time,
		})
	}
	return DB.Create(op_stat_slice).Error
}

func DelStats(id uint) error {
	return DB.Table(Table_Stats).Where("id = ?", id).Delete(&Stats{}).Error
}

func DelStatsByUserID(id uint) error {
	return DB.Table(Table_Stats).Where("user_id = ?", id).Delete(&Stats{}).Error
}

func UpdateStats(o *Stats) error {
	return DB.Save(o).Error
}

func UpdateStatsSlice(os []*Stats) error {
	return DB.Save(os).Error
}

func UpdateStatsByStats(osm map[TypeUserIDOpID]*Stats) error {
	// all user's opstats data
	t_c_os, e := GetAllStats()
	if e != nil {
		return e
	}
	// just for update
	c_os_for_update := make([]*Stats, 0)
	for _, c_os := range t_c_os {
		t_key := TypeUserIDOpID{
			UserId: c_os.UserID,
			OpId:   c_os.OpID,
		}
		// add new stats
		if _, ok := osm[t_key]; ok {
			t_v := osm[t_key]
			c_os.Success += t_v.Success
			c_os.Failure += t_v.Failure
			c_os.LastTime = t_v.LastTime
			c_os_for_update = append(c_os_for_update, c_os)
		}
	}
	return UpdateStatsSlice(c_os_for_update)
}

func GetAllStats() ([]*Stats, error) {
	ops := make([]*Stats, 0)
	e := DB.Table(Table_Stats).Find(&ops).Error
	return ops, e
}

func GetStatsByUserID(id uint) ([]*Stats, error) {
	ops := make([]*Stats, 0)
	e := DB.Table(Table_Stats).Where("user_id = ?", id).Find(&ops).Error
	return ops, e
}

func GetStatsByTgIDWithAlias(tg_id int64) ([]*ResultStatsForNotify, error) {
	result := make([]*ResultStatsForNotify, 0)
	select_str := fmt.Sprintf("select %v.alias as app_alias, %v.ms_username as user_alias, %v.alias as op_alias, %v.tg_id, %v.success, %v.failure, %v.last_time",
		Table_APPs, Table_Users, Table_Op, Table_Users, Table_Stats, Table_Stats, Table_Stats)
	select_str = fmt.Sprintf("%v from %v, %v, %v, %v", select_str, Table_APPs, Table_Users, Table_Op, Table_Stats)
	select_str = fmt.Sprintf("%v where %v.id = %v.app_id and %v.user_id = %v.id and %v.op_id = %v.id and %v.tg_id = %v",
		select_str, Table_APPs, Table_Users, Table_Stats, Table_Users, Table_Stats, Table_Op, Table_Users, tg_id)
	select_str = fmt.Sprintf("%s and %s.deleted_at IS NULL and %s.deleted_at IS NULL and %s.deleted_at IS NULL and %s.deleted_at IS NULL",
		select_str, Table_APPs, Table_Users, Table_Stats, Table_Op)
	e := DB.Raw(select_str).
		Scan(&result).Error
	return result, e
}

func GetAllStatsWithAlias() ([]*ResultStatsForNotify, error) {
	result := make([]*ResultStatsForNotify, 0)
	select_str := fmt.Sprintf("select %v.alias as app_alias, %v.ms_username as user_alias, %v.alias as op_alias, %v.tg_id, %v.success, %v.failure, %v.last_time",
		Table_APPs, Table_Users, Table_Op, Table_Users, Table_Stats, Table_Stats, Table_Stats)
	select_str = fmt.Sprintf("%v from %v, %v, %v, %v", select_str, Table_APPs, Table_Users, Table_Op, Table_Stats)
	select_str = fmt.Sprintf("%v where %v.id = %v.app_id and %v.user_id = %v.id and %v.op_id = %v.id",
		select_str, Table_APPs, Table_Users, Table_Stats, Table_Users, Table_Stats, Table_Op)
	select_str = fmt.Sprintf("%s and %s.deleted_at IS NULL and %s.deleted_at IS NULL and %s.deleted_at IS NULL and %s.deleted_at IS NULL",
		select_str, Table_APPs, Table_Users, Table_Stats, Table_Op)
	e := DB.Raw(select_str).
		Scan(&result).Error
	return result, e
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////
// table op_details
// ///////////////////////////////////////////////////////////////////////////////////////////////////////
// record every single operation details
type OpDetails struct {
	gorm.Model
	UserID    uint  `gorm:"column:user_id;not null"`
	OpID      uint  `gorm:"column:op_id;not null"`
	Success   int   `gorm:"column:success;not null"`
	Failure   int   `gorm:"column:failure;not null"`
	StartTime int64 `gorm:"column:start_time;not null"`
	EndTime   int64 `gorm:"column:end_time;not null"`
	Duration  int64 `gorm:"column:duration;not null"` // time in milliseconds
}

func (d *OpDetails) TableName() string {
	return Table_OpDetails
}

func SaveOpDetails(os []*OpDetails) error {
	return DB.Save(os).Error
}

func GetAllOpDetails() ([]*OpDetails, error) {
	ops := make([]*OpDetails, 0)
	err := DB.Table(Table_OpDetails).Find(&ops).Error
	return ops, err
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////
// table task_records
// ///////////////////////////////////////////////////////////////////////////////////////////////////////
// record details of every single task run in background
type TaskRecords struct {
	gorm.Model
	StartTime int64 `gorm:"column:start_time;not null"`
	EndTime   int64 `gorm:"column:end_time;not null"`
	Duration  int64 `gorm:"column:duration;not null"` // time in milliseconds
}

func (t *TaskRecords) TableName() string {
	return Table_TaskRecords
}

func SaveTaskRecords(ts []*TaskRecords) error {
	return DB.Save(ts).Error
}

func GetAllTaskRecords() ([]*TaskRecords, error) {
	ts := make([]*TaskRecords, 0)
	err := DB.Table(Table_TaskRecords).Find(&ts).Error
	return ts, err
}

// ////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////
