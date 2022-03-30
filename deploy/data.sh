# Load EVE dev resources
echo "download sde.zip from remote"
curl -O https://eve-static-data-export.s3-eu-west-1.amazonaws.com/tranquility/sde.zip

echo "unzip sde.zip"
unzip sde.zip >/dev/null 2>&1

num=$(ls | grep .yaml | wc -l)
if [ $num -gt 0 ]; then
        echo "cleanup the old yaml files"
        rm *.yaml
fi

echo "save the yaml files"
mv sde/fsd/typeIDs.yaml .
mv sde/fsd/factions.yaml .
mv sde/fsd/npcCorporations.yaml .

echo "loading the yaml files completed"
rm -rf sde
rm sde.zip

# Load Fuzzwork latest resources
num=$(ls | grep .csv | wc -l)
if [ $num -gt 0 ]; then
        echo "cleanup the old csv files"
        rm *.csv
fi

echo "download csv files from remote"
curl -O https://www.fuzzwork.co.uk/dump/latest/industryActivityProducts.csv
curl -O https://www.fuzzwork.co.uk/dump/latest/industryActivityMaterials.csv

echo "loading the csv files completed"

ls -lh