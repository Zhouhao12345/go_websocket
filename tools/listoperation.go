package tools

import "strconv"

func Sort(list []string) ([]string, error)  {
	length := len(list)
	for i:=0; i<length; i++ {
		for j:=i; j>0; j-- {
			now, err := strconv.Atoi(list[j])
			if err != nil {
				return make([]string, 1), err
			}
			next, err := strconv.Atoi(list[j-1])
			if err != nil {
				return make([]string, 1), err
			}
			if now < next {
				element := list[j]
				list[j] = list[j-1]
				list[j-1] = element
			}
		}
	}
	return list, nil
}
