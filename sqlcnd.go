package simple

import "github.com/jinzhu/gorm"

type SqlCnd struct {
	Params []ParamPair
	Orders []OrderByCol
	Paging *Paging
}

func NewSqlCnd(query string, args ...interface{}) *SqlCnd {
	q := &SqlCnd{}
	q.Where(query, args)
	return q
}

func (s *SqlCnd) Eq(column string, args ...interface{}) *SqlCnd {
	s.Where(column+" = ?", args)
	return s
}

func (s *SqlCnd) NotEq(column string, args ...interface{}) *SqlCnd {
	s.Where(column+" <> ?", args)
	return s
}

func (s *SqlCnd) Gt(column string, args ...interface{}) *SqlCnd {
	s.Where(column+" > ?", args)
	return s
}

func (s *SqlCnd) Gte(column string, args ...interface{}) *SqlCnd {
	s.Where(column+" >= ?", args)
	return s
}

func (s *SqlCnd) Lt(column string, args ...interface{}) *SqlCnd {
	s.Where(column+" < ?", args)
	return s
}

func (s *SqlCnd) Lte(column string, args ...interface{}) *SqlCnd {
	s.Where(column+" <= ?", args)
	return s
}

func (s *SqlCnd) Like(column string, str string) *SqlCnd {
	s.Where(column+" LIKE ?", "%"+str+"%")
	return s
}

func (s *SqlCnd) Where(query string, args ...interface{}) *SqlCnd {
	s.Params = append(s.Params, ParamPair{query, args})
	return s
}

func (s *SqlCnd) Asc(column string) *SqlCnd {
	s.Orders = append(s.Orders, OrderByCol{Column: column, Asc: true})
	return s
}

func (s *SqlCnd) Desc(column string) *SqlCnd {
	s.Orders = append(s.Orders, OrderByCol{Column: column, Asc: false})
	return s
}

func (s *SqlCnd) Limit(limit int) *SqlCnd {
	s.Page(1, limit)
	return s
}

func (s *SqlCnd) Page(page, limit int) *SqlCnd {
	if s.Paging == nil {
		s.Paging = &Paging{Page: page, Limit: limit}
	} else {
		s.Paging.Page = page
		s.Paging.Limit = limit
	}
	return s
}

func (s *SqlCnd) Query(db *gorm.DB) *gorm.DB {
	ret := db

	// where
	if len(s.Params) > 0 {
		for _, param := range s.Params {
			ret = ret.Where(param.Query, param.Args...)
		}
	}

	// order
	if len(s.Orders) > 0 {
		for _, order := range s.Orders {
			if order.Asc {
				ret = ret.Order(order.Column + " ASC")
			} else {
				ret = ret.Order(order.Column + " DESC")
			}
		}
	}

	// limit
	if s.Paging != nil && s.Paging.Limit > 0 {
		ret = ret.Limit(s.Paging.Limit)
	}

	// offset
	if s.Paging != nil && s.Paging.Offset() > 0 {
		ret = ret.Offset(s.Paging.Offset())
	}
	return ret
}

func (s *SqlCnd) Count(db *gorm.DB) *gorm.DB {
	ret := db

	// where
	if len(s.Params) > 0 {
		for _, query := range s.Params {
			ret = ret.Where(query.Query, query.Args...)
		}
	}

	return ret
}