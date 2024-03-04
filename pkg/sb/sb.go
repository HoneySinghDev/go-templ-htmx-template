package sb

import "github.com/nedpals/supabase-go"

func InitSB(supabaseURL string, supabaseKey string) *supabase.Client {
	sb := supabase.CreateClient(supabaseURL, supabaseKey)
	return sb
}
