package todoist

import "context"

type Section struct {
	HaveID
	HaveProjectID
	Collapsed    bool   `json:"collapsed"`
	Name         string `json:"name"`
	IsArchived   bool   `json:"is_archived"`
	IsDeleted    bool   `json:"is_deleted"`
	SectionOrder int    `json:"section_order"`
}

type Sections []Section

func (a Sections) GetIDByName(name string) string {
	for _, s := range a {
		if s.Name == name {
			return s.GetID()
		}
	}
	return ""
}

func (section Section) AddParam() interface{} {
	param := map[string]interface{}{}
	if section.Name != "" {
		param["name"] = section.Name
	}
	if section.ProjectID != "" {
		param["project_id"] = section.ProjectID
	}
	return param
}

func (c *Client) AddSection(ctx context.Context, section Section) error {
	commands := Commands{
		NewCommand("section_add", section.AddParam()),
	}
	return c.ExecCommands(ctx, commands)
}

func (c *Client) DeleteSection(ctx context.Context, id string) error {
	commands := Commands{
		NewCommand("section_delete", map[string]interface{}{"id": id}),
	}
	return c.ExecCommands(ctx, commands)
}

func (c *Client) ArchiveSection(ctx context.Context, id string) error {
	commands := Commands{
		NewCommand("section_archive", map[string]interface{}{"id": id}),
	}
	return c.ExecCommands(ctx, commands)
}

func (c *Client) UnarchiveSection(ctx context.Context, id string) error {
	commands := Commands{
		NewCommand("section_unarchive", map[string]interface{}{"id": id}),
	}
	return c.ExecCommands(ctx, commands)
}

func (c *Client) UpdateSection(ctx context.Context, id string, name string) error {
	commands := Commands{
		NewCommand("section_update", map[string]interface{}{"id": id, "name": name}),
	}
	return c.ExecCommands(ctx, commands)
}

func (c *Client) MoveSection(ctx context.Context, sectionID string, projectID string) error {
	commands := Commands{
		NewCommand("section_move", map[string]interface{}{"id": sectionID, "project_id": projectID}),
	}
	return c.ExecCommands(ctx, commands)
}

func (c *Client) ReorderSections(ctx context.Context, ids []string) error {
	type sectionOrder struct {
		ID           string `json:"id"`
		SectionOrder int    `json:"section_order"`
	}
	sections := make([]sectionOrder, len(ids))
	for i, id := range ids {
		sections[i] = sectionOrder{ID: id, SectionOrder: i + 1}
	}
	commands := Commands{
		NewCommand("section_reorder", map[string]interface{}{"sections": sections}),
	}
	return c.ExecCommands(ctx, commands)
}
