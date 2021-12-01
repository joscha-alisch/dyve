export const FormTeamGeneral = {
    schema: {
        name: {
            type: "string",
            title: "Team Name",
            hint: "The name of your team",
            required: true,
        },
        slug: {
            type: "string",
            title: "Team Slug",
            validators: ["noSpaces"],
            hint: "A concise identifier for URLs and internals",
            required: true
        },
        description: {
            type: "string",
            hint: "Describe what your team does",
            title: "Team Description",
            multiline: true,
        },
    },
    ui: [
        {
            layout: "row",
            children: [
                {
                    layout: "column",
                    children: ["name", "slug"]
                },
                "description",
            ],
        },
    ]
}

export const FormTeamAccess = {
    schema: {
        admin: {
            type: "array",
            title: "Admin Groups",
            hint: "Configure the groups that should be admins in this team",
            required: true,
        },
        member: {
            type: "array",
            title: "Member Groups",
            hint: "Configure the groups that should be team members",
            required: true,
        },
        viewer: {
            type: "array",
            title: "Viewer Groups",
            hint: "Configure the groups that should be able to view minimal info about this teams resources",
            required: true,
        },
    },
    ui: [
        {
            layout: "column",
            children: [
                {
                    field: "admin",
                    component: "tags",
                    dataKey: "groups",
                    labelKey: "name",
                    groupKey: "providerName",
                    style: {
                        minHeight: "100px"
                    }
                },
                {
                    field: "member",
                    component: "tags",
                    dataKey: "groups",
                    labelKey: "name",
                    groupKey: "providerName",
                    style: {
                        minHeight: "100px"
                    }
                },
                {
                    field: "viewer",
                    component: "tags",
                    dataKey: "groups",
                    labelKey: "name",
                    groupKey: "providerName",
                    style: {
                        minHeight: "100px"
                    }
                }
            ]
        }
    ]
}

export const FormTeamAssociations = {
    schema: {
        apps: {
            type: "array",
            title: "App Associations",
            hint: "Apps with these labels belong to this team",
            required: true,
        },
        pipelines: {
            type: "array",
            title: "Pipeline Associations",
            hint: "Pipelines with these labels belong to this team",
            required: true,
        },
    },
    ui: [
        {
            layout: "row",
            children: [
                {
                    field: "apps",
                    component: "conditions",
                },
                {
                    field: "pipelines",
                    component: "conditions",
                }
            ]
        }

    ],
}