import {Typography} from "@mui/material";

export default function HomePage() {
    return (
        <>
            <Typography variant="h4">Play!</Typography>
        </>
    );
}


<form action="/create-room" method="post">
    <br/>
    <p>Language:</p>
    {{range .languages}}
    {{if eq. "en"}}
    <input type="radio" id="{{ . }}" name="language" value="{{ . }}" required="required" checked/>
    {{else}}
    <input type="radio" id="{{ . }}" name="language" value="{{ . }}" required="required"/>
    {{end}}
    <label htmlFor="{{ . }}">{{.}}</label>
    <br/>
    {{end}}
    <label htmlFor="allow-rude-words">Allow rude words</label>
    <input type="checkbox" id="allow-rude-words" name="allow-rude-words"/>
    <br/>
    <label htmlFor="only-external-vocabulary">Use only additional vocabulary</label>
    <input type="checkbox" id="only-external-vocabulary" name="only-external-vocabulary"/>
    <br/>
    <label htmlFor="additional-vocabulary">Additional vocabulary</label>
    <textarea id="additional-vocabulary" name="additional-vocabulary" placeholder="word1,word2,word3" rows="3"
              style="width: 200px;"></textarea>
    <br/>
    <label htmlFor="clock">Clock in seconds. -1 for no clock.</label>
    <input type="number" id="clock" name="clock" value="60"/>
    <br/>
    <input type="submit"/>
</form>
