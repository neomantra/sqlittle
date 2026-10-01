package sqlittle

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestColumns(t *testing.T) {
	db, err := Open("testdata/words.sqlite")
	require.NoError(t, err)
	defer db.Close()

	t.Run("fine", func(t *testing.T) {
		cols, err := db.Columns("words")
		require.NoError(t, err)
		require.Equal(t, []string{"word", "length"}, cols)
	})

	t.Run("no such table", func(t *testing.T) {
		cols, err := db.Columns("notwords")
		require.EqualError(t, err, `no such table: "notwords"`)
		require.Nil(t, cols)
	})
}

func TestOpenBytes(t *testing.T) {
	b, err := os.ReadFile("testdata/words.sqlite")
	require.NoError(t, err)

	words := func(db *DB) []string {
		var res []string
		require.NoError(t, db.Select("words", func(r Row) {
			var w string
			require.NoError(t, r.Scan(&w))
			res = append(res, w)
		}, "word"))
		return res
	}

	fromFile, err := Open("testdata/words.sqlite")
	require.NoError(t, err)
	defer fromFile.Close()
	fromBytes, err := OpenBytes(b)
	require.NoError(t, err)
	defer fromBytes.Close()

	want := words(fromFile)
	require.NotEmpty(t, want)
	require.Equal(t, want, words(fromBytes))

	_, err = OpenBytes(nil)
	require.Error(t, err)
}
