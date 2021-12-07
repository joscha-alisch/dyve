
compPath="src/components/$1"
upperName=$2
lowerName="$(tr '[:upper:]' '[:lower:]' <<< "${upperName:0:1}")${upperName:1}"

mkdir -p "$compPath/$lowerName"
for i in ./scripts/template/*; do
    fName=$(basename "$i")
    newFile="$compPath/$lowerName/$lowerName.${fName#*.}"
    cp "$i" "$newFile"

    sed -i '' "s/template_lower/$lowerName/g" "$newFile"
    sed -i '' "s/template_upper/$upperName/g" "$newFile"
done