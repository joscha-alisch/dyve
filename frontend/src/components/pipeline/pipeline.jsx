import CytoscapeComponent from "react-cytoscapejs/src/component";
import styles from "./pipeline.module.sass"

const Pipeline = () => {
    const elements = [
        { data: { id: 'one', label: 'Node 1' }, position: { x: 50, y: 50 } },
        { data: { id: 'two', label: 'Node 2' }, position: { x: 200, y: 200 } },
        { data: { source: 'one', target: 'two', label: 'Edge from Node1 to Node2' } }
    ];

    return <CytoscapeComponent
                panningEnabled={false}
                zoomingEnabled={false}
                autoungrabify={true}
                autolock={true}
                autounselectify={true}
                elements={elements}
                className={styles.Graph}
                stylesheet={[
                    {
                        selector: 'node',
                        style: {
                            width: 50,
                            height: 100,
                            shape: 'rectangle'
                        }
                    },
                    {
                        selector: 'edge',
                        style: {
                            width: 1
                        }
                    }
                ]}/>
}

export default Pipeline