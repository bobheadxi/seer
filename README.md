<p align="center">
  <img src="https://vignette.wikia.nocookie.net/leagueoflegends/images/3/3c/Seer_Stone_item.png/revision/latest?cb=20171221231955" />
</p>

<br />

<p align="center">
  <a href="https://godoc.org/github.com/bobheadxi/seer">
    <img src="https://godoc.org/github.com/bobheadxi/seer?status.svg" alt="GoDoc">
  </a>

  <a href="https://dev.azure.com/bobheadxi/bobheadxi/_build/latest?definitionId=8&branchName=master">
    <img src="https://dev.azure.com/bobheadxi/bobheadxi/_apis/build/status/bobheadxi.seer?branchName=master"
      alt="CI Status" />
  </a>

  <a href="https://seer-engine.herokuapp.com/status">
    <img src="https://img.shields.io/website/https/seer-engine.herokuapp.com/status.svg?down_color=lightgrey&down_message=offline&label=api&up_message=online"
      alt="API Status" >
  </a>

  <a href="https://seer.bobheadxi.dev">
    <img src="https://img.shields.io/website/https/seer.bobheadxi.dev.svg?down_color=lightgrey&down_message=offline&up_message=online"
      alt="Website Status">
  </a>
</p>

<br />

## What is this?

Seer aims to be a service where users can input a group of [League of Legends](https://na.leagueoflegends.com/en/)
players to track as a "team". The service will then pull game matches and
associated data periodically or on trigger from the [Riot Games API](https://developer.riotgames.com/)
and pipe them into a public [BigQuery](https://cloud.google.com/bigquery/)
dataset. The goal is to make this dataset available for everyone to play around
with, and for it to power the main [seer.bobheadxi.dev](https://seer.bobheadxi.dev)
website. This website aims to allow members of teams to gain insights into their
team's performance over time, as well as to allow players to "scout" other teams.
