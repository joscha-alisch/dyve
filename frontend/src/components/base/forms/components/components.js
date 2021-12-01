import {withComponents} from "../../../../packages/formidable/components";
import {associationsField} from "./associations";
import {arrayField} from "./array";
import {conditionsField} from "./conditions";
import {tagsField} from "./tags";

export const components = withComponents({
    array: {
        default: arrayField,
        associations: associationsField,
        conditions: conditionsField,
        tags: tagsField
    }
})