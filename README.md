

var conference Conference
var conferences []Conference

db.Select(&conference, dal.Where("name=?", "bacon"))
db.Select(&conferences, dal.Where("closes < now()").Limit(50))
db.Insert(&conference)
