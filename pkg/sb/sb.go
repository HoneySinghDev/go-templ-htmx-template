package sb

import "github.com/nedpals/supabase-go"

func InitSB(supabaseUrl string, supabaseKey string) *supabase.Client {
	sb := supabase.CreateClient(supabaseUrl, supabaseKey)
	return sb
}
