import React, {useEffect, useState} from "react"
import styles from "./multiPageForm.module.sass"
import PropTypes from "prop-types"
import Box from "../../box/box";
import Tabs, {Tab} from "../../tabs/tabs";
import {Form} from "../../../../packages/formidable";
import {components} from "../components/components";
import {Button} from "@mui/material";
import {useNotifications} from "../../../../context/notifications";

const MultiPageForm = ({className, forms, getData, submitText = "Submit", onSubmit = (data, notify) => {}}) => {
    let [state, setState] = useState({})
    let [data, setData] = useState({})
    let {notify} = useNotifications()

    let updateForm = (form) => (formState) => {
        setState({
            ...state,
            [form]: formState
        })
    }

    useEffect(() => {
        async function fetchData() {
            if (getData) {
                let data = await getData()
                setData(data)
            }
        }

        fetchData()
    }, [getData])

    let buttons = (index, setIndex) => {
        let nextFunc = () => setIndex((index + 1))
        if (index + 1 === forms.length) {
            nextFunc = () => {
                let submitData = {}
                for (let i = 0; i < forms.length; i++) {
                    if (forms[i].mappingKey) {
                        submitData = {
                            ...submitData,
                            [forms[i].mappingKey]: state[i]
                        }
                    } else {
                        submitData = {
                            ...submitData,
                            ...state[i]
                        }
                    }

                }
                onSubmit(submitData, notify)
            }
        }

        return <nav className={styles.Buttons}>
            {index > 0 && <Button variant={"outlined"} onClick={() => setIndex((index - 1))}>Previous</Button>}
            <Button onClick={nextFunc} variant={"contained"}>{index + 1 === forms.length ? submitText : "Next"}</Button>
        </nav>
    }

    return <Tabs renderHeader={buttons} renderFooter={buttons}>
        {forms.map(({title, icon, form}, i) => <Tab key={i} title={title}>
            <Box className={styles.Box}>
                <Form {...form} components={components} data={data} state={state[i]} setState={updateForm(i)}/>
            </Box>
        </Tab>)}
    </Tabs>
}

MultiPageForm.propTypes = {
    className: PropTypes.string,
    forms: PropTypes.arrayOf(PropTypes.shape({
        title: PropTypes.string,
        icon: PropTypes.object,
        form: PropTypes.object,
        mappingKey: PropTypes.string,
    }))
}

export default MultiPageForm