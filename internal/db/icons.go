package db

type Icon struct {
	ID      int    `db:"iconid"`
	URL     string `db:"iconurl"`
	ToolTip string `db:"tooltip"`
	RoleID  int    `db:"roleid"`
	Paths   []string
}

func (i *Icon) getPaths(e *EpilogueConnection) error {
	var paths []string
	err := e.DB.Select(&paths, "SELECT iconpath FROM iconpaths WHERE iconid=?", i.ID)
	if err != nil {
		return err
	}
	i.Paths = paths
	return nil
}

func (e *EpilogueConnection) GetIconsByRoleID(role int) ([]*Icon, error) {
	var (
		icons = []*Icon{}
	)
	err := e.DB.Select(&icons, "SELECT * FROM icons WHERE roleid<=?", role)
	if err != nil {
		return icons, err
	}

	for _, icon := range icons {
		err = icon.getPaths(e)
		if err != nil {
			return icons, err
		}
	}
	return icons, nil
}
