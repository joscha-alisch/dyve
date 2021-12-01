import React from "react"
import MultiPageForm from "./multiPageForm";

export default {
    title: 'Components/Forms/Multi Page',
    component: MultiPageForm,
}

export const StoryMultiStepForm = (args) => <MultiPageForm {...args}/>

StoryMultiStepForm.storyName = "Multi Page"
StoryMultiStepForm.args = {
    forms: [
        {
            title: "Form 1",
            form: {
                data: {
                    someProperty: {
                        type: "string"
                    }
                }
            }
        },
        {
            title: "Form 2",
            form: {
                data: {
                    someOtherProperty: {
                        type: "string"
                    }
                }
            }
        }
    ]
}