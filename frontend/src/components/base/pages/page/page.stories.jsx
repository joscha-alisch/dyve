import React from "react"
import Page from "./page";

export default {
    title: 'App/Pages/Page',
    component: Page,
}

export const StoryPage = (args) => <Page {...args}/>

StoryPage.storyName = "Page"
StoryPage.args = {
    title: "Current Page",
    parent: "Parent Page",
    children: <React.Fragment>
        <p>
            Lorem ipsum dolor sit amet, consectetur adipisicing elit. Aperiam doloribus, fugiat maxime officiis qui quis
            sunt. Adipisci at commodi consectetur impedit laudantium obcaecati optio qui repellat sapiente, tenetur
            voluptate voluptatibus. Lorem ipsum dolor sit amet, consectetur adipisicing elit. Accusantium, aliquam
            dolorem enim et expedita facilis hic impedit maiores maxime neque nesciunt obcaecati quasi quos
            reprehenderit sequi similique tempora velit voluptatum?
        </p>
    </React.Fragment>
}