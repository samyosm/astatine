use regex::Regex;

pub fn wordpunkt(text: &str) -> Vec<&str> {
    Regex::new(r"\w+|[^\w\s]+")
        .unwrap()
        .find_iter(text)
        .map(|g| g.as_str())
        .collect()
}
