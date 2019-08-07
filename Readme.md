Создать проект
```
Method: PUT
Resource: /project/
Body: 
{
    title: string
}
Response: project key
```

Добавить новую запись с метриками:

```
Method: PUT
Resource: /log/${project_key}/
Body: 
{
    metricName:     string,
    durationMs:     string,
    additionalInfo: object,
}
```