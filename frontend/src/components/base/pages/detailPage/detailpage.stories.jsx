import React from "react"
import DetailPage from "./detailpage";

export default {
    title: 'Components/DetailPage',
    component: DetailPage,
}

export const StoryDetailPage= (args) => <DetailPage {...args}/>

StoryDetailPage.storyName = "DetailPage"
StoryDetailPage.args = {
}