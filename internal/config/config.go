package config

type File interface {
	GetKey(name string) string
	RemoveKey(name string)
	SetKey(name, value string)
	Save() error
}

type Config struct {
	// Dirty Flags
	systemDirty  bool
	userDirty    bool
	projectDirty bool
	// Config Files
	System  File
	User    File
	Project File
}

func (c *Config) Save() error {
	if c.projectDirty {
		if err := c.Project.Save(); err != nil {
			return err
		}
	}
	if c.userDirty {
		if err := c.User.Save(); err != nil {
			return err
		}
	}
	if c.systemDirty {
		if err := c.System.Save(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) GetPinned() string {
	return c.Project.GetKey("core.version")
}

func (c *Config) SetPinned(version string) (before, after string) {
	before = c.Project.GetKey("core.version")
	if before == version {
		after = "no change"
	} else {
		if version == "" {
			after = "not pinned"
			c.Project.RemoveKey("core.version")
		} else {
			after = version
			c.Project.SetKey("core.version", version)
		}
		c.projectDirty = true
	}
	if before == "" {
		before = "not pinned"
	}
	return
}
