package params

import "fmt"

type Sort struct {
	Id   int64 `json:"id" binding:"required"`
	Sort int   `json:"sort" binding:"omitempty,max=1000000000"`
}

type BatchSort struct {
	Sorts []*Sort `json:"sorts" binding:"required,min=1,max=200,dive"`
}

func (b *BatchSort) Verify() error {
	length := len(b.Sorts)
	tempMap := make(map[int64]int, length)
	for i := 0; i < length; i++ {
		sort := b.Sorts[i]
		tempMap[sort.Id] = tempMap[sort.Id] + 1
	}
	// 验证重复的数据
	for k, v := range tempMap {
		if v > 1 {
			return fmt.Errorf("排序列表出现了重复 主键:%d 重复次数:%d", k, v)
		}
	}
	return nil
}

type SortMongo struct {
	IdParam
	Sort int `json:"sort" binding:"omitempty,max=1000000000"`
}

type BatchSortMongo struct {
	Sorts []*SortMongo `json:"sorts" binding:"required,min=1,max=200,dive"`
}

func (b *BatchSortMongo) Verify() error {
	length := len(b.Sorts)
	tempMap := make(map[string]int, length)
	for i := 0; i < length; i++ {
		sort := b.Sorts[i]
		tempMap[sort.Id] = tempMap[sort.Id] + 1
	}
	// 验证重复的数据
	for k, v := range tempMap {
		if v > 1 {
			return fmt.Errorf("排序列表出现了重复 主键:%s 重复次数:%d", k, v)
		}
	}
	return nil
}
