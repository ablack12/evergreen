ReactDOM.render(
  <Root
    project={window.project}
    userTz={window.userTz}
    palette={window.useAlternatePalette}
    jiraHost={window.jiraHost}
    user={window.user}
  />, document.getElementById('root')
);
